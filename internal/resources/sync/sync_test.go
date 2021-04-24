package sync

import (
	"testing"
	"time"

	"github.com/jinzhu/copier"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	fp "path/filepath"

	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/iinitialize"
	"github.com/kick-project/kick/internal/settings/isync"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
)

func setup() (*settings.Settings, *gorm.DB) {
	home := fp.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)
	init := initialize.Initialize{}
	err := copier.Copy(&init, iinitialize.Inject(s))
	if err != nil {
		errutils.Epanic(err)
	}
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

	inserted := false
	var err error
	for i := 0; i < 10; i++ {
		result := db.Clauses(clause.Insert{Modifier: "OR IGNORE"}).Create(&g)

		// TODO: Find internal race condition within gorm or sqlite3 library.
		if result.Error == nil {
			inserted = true
			break
		} else {
			err = result.Error
			time.Sleep(time.Duration(i) * 100 * time.Millisecond)
		}
	}
	if !inserted {
		t.Errorf("Could not insert into database: %v\n", err)
	}

	syncobj := &Sync{}
	err = copier.Copy(syncobj, isync.Inject(s))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	syncobj.Global()
}

func TestGlobalNoURL(t *testing.T) {
}
