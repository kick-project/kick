package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/utils"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
	"github.com/stretchr/testify/assert"
)

var driver string = "sqlite3"
var dsn string = "file:%s?_foreign_key=on"

func TestCreate(t *testing.T) {
	var err error
	db, path := tmpDb()
	_ = path
	CreateSchema(db)
	defer db.Close()
	stmt, err := db.Prepare(`SELECT count(*) as count FROM master WHERE url="none"`)
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var count int
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, count)
}

func TestCreateModel(t *testing.T) {
	path := filepath.Join(utils.TempDir(), "schema_model_test.db")
	_, err := os.Stat(path)
	if err == nil {
		os.Remove(path)
	}

	CreateModel(&ModelOptions{
		File: path,
	})
}

// tmpDb create a temporary db file
// the cleanup function should be called to remove the database after exict
func tmpDb() (db *sql.DB, path string) {
	path = filepath.Join(utils.TempDir(), "schema_test.db")
	_, err := os.Stat(path)
	if err == nil {
		os.Remove(path)
	}
	db, err = sql.Open(driver, fmt.Sprintf(dsn, path))
	if err != nil {
		panic(err)
	}
	return db, path
}
