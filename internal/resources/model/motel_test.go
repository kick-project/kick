package model_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/utils"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
)

func TestCreateModel(t *testing.T) {
	path := filepath.Join(utils.TempDir(), "model_test.db")
	_, err := os.Stat(path)
	if err == nil {
		os.Remove(path)
	}

	model.CreateModel(&model.Options{
		File: path,
	})
}
