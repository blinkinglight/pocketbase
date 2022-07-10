//go:build pq

package core

import (
	_ "github.com/lib/pq"
	"github.com/pocketbase/dbx"
	"os"
	"strings"
)

func connectDB(dbPath string) (*dbx.DB, error) {
	if strings.Contains(dbPath, "logs.db") {
		return dbx.MustOpen("postgres", os.Getenv("LOGS_DATABASE"))
	}
	return dbx.MustOpen("postgres", os.Getenv("DATABASE"))
}
