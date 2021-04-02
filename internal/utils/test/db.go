package test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'

	"github.com/kick-project/kick/internal/utils"
)

var driver string = "sqlite3"
var dsn string = "file:%s?_foreign_key=on"

// TmpDb create a temporary db file
// the cleanup function should be called to remove the database after exict
func TmpDb(t *testing.T) (db *sql.DB, path string) {
	path = filepath.Join(utils.TempDir(), "sync_test.db")
	_, err := os.Stat(path)
	if err == nil {
		os.Remove(path)
	}
	db, err = sql.Open(driver, fmt.Sprintf(dsn, path))
	if err != nil {
		t.Errorf("Error %v", err)
	}
	return db, path
}
