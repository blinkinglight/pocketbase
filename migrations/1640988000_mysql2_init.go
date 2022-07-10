//go:build mysql

// Package migrations contains the system PocketBase DB migrations.
package migrations

import (
	"fmt"

	"github.com/blinkinglight/pocketbase-mysql/daos"
	"github.com/blinkinglight/pocketbase-mysql/models"
	"github.com/blinkinglight/pocketbase-mysql/models/schema"
	"github.com/blinkinglight/pocketbase-mysql/tools/migrate"
	"github.com/pocketbase/dbx"
)

var AppMigrations migrate.MigrationsList

// Register is a short alias for `AppMigrations.Register()`
// that is usually used in external/user defined migrations.
func Register(
	up func(db dbx.Builder) error,
	down func(db dbx.Builder) error,
	optFilename ...string,
) {
	AppMigrations.Register(up, down, optFilename...)
}

func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		_, tablesErr := db.NewQuery(`
			CREATE TABLE if not exists _admins (
				id              TEXT NOT NULL,
				avatar          INTEGER DEFAULT 0 NOT NULL,
				email           TEXT NOT NULL DEFAULT '' ,
				tokenKey        TEXT NOT NULL DEFAULT '' ,
				passwordHash    TEXT NOT NULL,
				lastResetSentAt TEXT NOT NULL DEFAULT '' ,
				created         TEXT NOT NULL DEFAULT '' ,
				updated         TEXT NOT NULL DEFAULT '' 
			) ENGINE = InnoDB;
			`).Execute()
		if tablesErr != nil {
			return tablesErr
		}
		_, tablesErr = db.NewQuery(`
			CREATE TABLE if not exists _users (
				id                     TEXT NOT NULL,
				verified               BOOLEAN DEFAULT FALSE NOT NULL,
				email                  TEXT NOT NULL,
				tokenKey               TEXT NOT NULL,
				passwordHash           TEXT NOT NULL,
				lastResetSentAt        TEXT NOT NULL DEFAULT '' ,
				lastVerificationSentAt TEXT NOT NULL DEFAULT '' ,
				created                TEXT NOT NULL DEFAULT '' ,
				updated                TEXT NOT NULL DEFAULT '' 
			) ENGINE = InnoDB;
			`).Execute()
		if tablesErr != nil {
			return tablesErr
		}
		_, tablesErr = db.NewQuery(`
			CREATE TABLE if not exists _collections (
				id         TEXT NOT NULL,
				system     BOOLEAN NOT NULL DEFAULT FALSE,
				name       TEXT NOT NULL,
				` + "`" + `schema` + "`" + `JSON NOT NULL DEFAULT '{}',
				listRule   TEXT NULL,
				viewRule   TEXT NULL,
				createRule TEXT NULL,
				updateRule TEXT NULL,
				deleteRule TEXT NULL,
				created    TEXT NOT NULL DEFAULT '' ,
				updated    TEXT NOT NULL DEFAULT '' 
			) ENGINE = InnoDB;
			`).Execute()
		if tablesErr != nil {
			return tablesErr
		}
		_, tablesErr = db.NewQuery(`
			CREATE TABLE if not exists _params (
				id      TEXT NOT NULL,
				` + "`" + `key` + "`" + `     TEXT NOT NULL,
				` + "`" + `value` + "`" + `   JSON DEFAULT NULL,
				created TEXT NOT NULL DEFAULT '' ,
				updated TEXT NOT NULL DEFAULT '' 
			) ENGINE = InnoDB;
		`).Execute()
		if tablesErr != nil {
			return tablesErr
		}

		// inserts the system profiles collection
		// -----------------------------------------------------------

		profileOwnerRule := fmt.Sprintf("%s = @request.user.id", models.ProfileCollectionUserFieldName)
		collection := &models.Collection{
			Name:       models.ProfileCollectionName,
			System:     true,
			CreateRule: &profileOwnerRule,
			ListRule:   &profileOwnerRule,
			ViewRule:   &profileOwnerRule,
			UpdateRule: &profileOwnerRule,
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     models.ProfileCollectionUserFieldName,
					Type:     schema.FieldTypeUser,
					Unique:   true,
					Required: true,
					System:   true,
					Options: &schema.UserOptions{
						MaxSelect:     1,
						CascadeDelete: true,
					},
				},
				&schema.SchemaField{
					Name:    "name",
					Type:    schema.FieldTypeText,
					Options: &schema.TextOptions{},
				},
				&schema.SchemaField{
					Name: "avatar",
					Type: schema.FieldTypeFile,
					Options: &schema.FileOptions{
						MaxSelect: 1,
						MaxSize:   5242880,
						MimeTypes: []string{
							"image/jpg",
							"image/jpeg",
							"image/png",
							"image/svg+xml",
							"image/gif",
						},
					},
				},
			),
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		tables := []string{
			"_params",
			"_collections",
			"_users",
			"_admins",
			models.ProfileCollectionName,
		}

		for _, name := range tables {
			if _, err := db.DropTable(name).Execute(); err != nil {
				return err
			}
		}

		return nil
	})
}
