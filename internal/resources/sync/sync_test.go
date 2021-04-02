package sync

import (
	"os"
	"testing"

	"github.com/jinzhu/copier"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	fp "path/filepath"

	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/isync"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
)

func setup() (*settings.Settings, *gorm.DB) {
	home := fp.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)
	err := os.MkdirAll(s.PathMetadataDir, 0755)
	if err != nil {
		errutils.Epanic(err)
	}
	model.CreateModel(&model.Options{
		File: s.ModelDB,
	})
	db := s.GetORM()
	return s, db
}

func TestGlobal(t *testing.T) {
	s, db := setup()

	// TODO: Create a global URL
	g := model.Global{
		Name: "master",
		URL:  "http://127.0.0.1:5000/master2.git",
		Desc: "master 2",
	}

	result := db.Clauses(clause.Insert{Modifier: "OR IGNORE"}).Create(&g)
	errutils.Efatal(result.Error)

	syncobj := &Sync{}
	err := copier.Copy(syncobj, isync.Inject(s))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	syncobj.Global()
}

func TestGlobalNoURL(t *testing.T) {
}
