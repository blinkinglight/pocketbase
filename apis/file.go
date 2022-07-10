package apis

import (
	"github.com/blinkinglight/pocketbase-mysql/core"
	"github.com/blinkinglight/pocketbase-mysql/models"
	"github.com/blinkinglight/pocketbase-mysql/models/schema"
	"github.com/blinkinglight/pocketbase-mysql/tools/list"
	"github.com/blinkinglight/pocketbase-mysql/tools/rest"
	"github.com/labstack/echo/v5"
)

var imageContentTypes = []string{"image/png", "image/jpg", "image/jpeg"}
var defaultThumbSizes = []string{"100x100"}

// BindFileApi registers the file api endpoints and the corresponding handlers.
func BindFileApi(app core.App, rg *echo.Group) {
	api := fileApi{app: app}

	subGroup := rg.Group("/files", ActivityLogger(app))
	subGroup.GET("/:collection/:recordId/:filename", api.download, LoadCollectionContext(api.app))
}

type fileApi struct {
	app core.App
}

func (api *fileApi) download(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return rest.NewNotFoundError("", nil)
	}

	recordId := c.PathParam("recordId")
	if recordId == "" {
		return rest.NewNotFoundError("", nil)
	}

	record, err := api.app.Dao().FindRecordById(collection, recordId, nil)
	if err != nil {
		return rest.NewNotFoundError("", err)
	}

	filename := c.PathParam("filename")

	fileField := record.FindFileFieldByFile(filename)
	if fileField == nil {
		return rest.NewNotFoundError("", nil)
	}
	options, _ := fileField.Options.(*schema.FileOptions)

	fs, err := api.app.NewFilesystem()
	if err != nil {
		return rest.NewBadRequestError("Filesystem initialization failure.", err)
	}
	defer fs.Close()

	originalPath := record.BaseFilesPath() + "/" + filename
	servedPath := originalPath
	servedName := filename

	// check for valid thumb size param
	thumbSize := c.QueryParam("thumb")
	if thumbSize != "" && (list.ExistInSlice(thumbSize, defaultThumbSizes) || list.ExistInSlice(thumbSize, options.Thumbs)) {
		// extract the original file meta attributes and check it existence
		oAttrs, oAttrsErr := fs.Attributes(originalPath)
		if oAttrsErr != nil {
			return rest.NewNotFoundError("", err)
		}

		// check if it is an image
		if list.ExistInSlice(oAttrs.ContentType, imageContentTypes) {
			// add thumb size as file suffix
			servedName = thumbSize + "_" + filename
			servedPath = record.BaseFilesPath() + "/thumbs_" + filename + "/" + servedName

			// check if the thumb exists:
			// - if doesn't exist - create a new thumb with the specified thumb size
			// - if exists - compare last modified dates to determine whether the thumb should be recreated
			tAttrs, tAttrsErr := fs.Attributes(servedPath)
			if tAttrsErr != nil || oAttrs.ModTime.After(tAttrs.ModTime) {
				if err := fs.CreateThumb(originalPath, servedPath, thumbSize, false); err != nil {
					servedPath = originalPath // fallback to the original
				}
			}
		}
	}

	event := &core.FileDownloadEvent{
		HttpContext: c,
		Record:      record,
		Collection:  collection,
		FileField:   fileField,
		ServedPath:  servedPath,
		ServedName:  servedName,
	}

	return api.app.OnFileDownloadRequest().Trigger(event, func(e *core.FileDownloadEvent) error {
		if err := fs.Serve(e.HttpContext.Response(), e.ServedPath, e.ServedName); err != nil {
			return rest.NewNotFoundError("", err)
		}

		return nil
	})
}
