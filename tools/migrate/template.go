package migrate

const createTemplateContent = `package migrations

import (
	"github.com/pocketbase/dbx"
	m "github.com/blinkinglight/pocketbase-mysql/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		// add up queries...

		return nil
	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
`
