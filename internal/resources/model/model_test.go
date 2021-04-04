package model

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestModel(t *testing.T) {
	p := filepath.Clean(fmt.Sprintf("%s/TestModel.db", utils.TempDir()))
	db, err := gorm.Open(sqlite.Open(p), &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		t.Fatalf("Can not create a ORM connection %s: %v", p, err)
	}
	for _, m := range []interface{}{&Global{}, &Master{}, &Template{}, &Installed{}, &Sync{}, &Versions{}} {
		err := db.AutoMigrate(m)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}
