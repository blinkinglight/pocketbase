//go:build mysql

package logs

import (
	"github.com/blinkinglight/pocketbase/tools/migrate"
	"github.com/pocketbase/dbx"
)

var LogsMigrations migrate.MigrationsList

func init() {
	LogsMigrations.Register(func(db dbx.Builder) (err error) {
		_, err = db.NewQuery(`

		CREATE TABLE {{_requests}} (
			id TEXT NOT NULL ,
			url TEXT NOT NULL DEFAULT '' ,
			method TEXT NOT NULL DEFAULT 'get' , 
			status INT NOT NULL DEFAULT '200' ,
			auth TEXT NOT NULL DEFAULT 'guest' ,
			ip TEXT NOT NULL DEFAULT '127.0.0.1' ,
			referer TEXT NOT NULL DEFAULT '' , 
			userAgent TEXT NOT NULL DEFAULT 'guest' , 
			meta JSON NOT NULL DEFAULT '{}',
			created TEXT NOT NULL DEFAULT '' ,
			updated TEXT NOT NULL DEFAULT '' ,
		) ENGINE = InnoDB;

		ALTER TABLE {{_requests}}
		ADD PRIMARY KEY (` + "`" + `id` + "`" + `(200));
		COMMIT;

		-- CREATE INDEX _request_status_idx on {{_requests}} (status);
		-- CREATE INDEX _request_auth_idx on {{_requests}} (auth);
		-- CREATE INDEX _request_ip_idx on {{_requests}} (ip);
		`).Execute()

		return err
	}, func(db dbx.Builder) error {
		_, err := db.DropTable("_requests").Execute()
		return err
	})
}
