package main

import (
	"log"

	"github.com/blinkinglight/pocketbase"
	"github.com/blinkinglight/pocketbase/apis"
	"github.com/blinkinglight/pocketbase/core"
	"github.com/labstack/echo/v5"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		subFs := echo.MustSubFS(e.Router.Filesystem, "pb_public")
		e.Router.GET("/*", apis.StaticDirectoryHandler(subFs, false))

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
