package sync_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"path/filepath"
	fp "path/filepath"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/resources/model/clauses"
	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/resources/testtools"
)

func setup(t *testing.T, home string, models ...interface{}) (*sync.Sync, *di.DI, *gorm.DB) {
	home = fp.Join(testtools.TempDir(), home)
	inject := di.Setup(home)

	init := inject.MakeInitialize()
	init.Init()

	db := inject.MakeORM()

	for _, m := range models {
		inserted := false
		var err error
		for i := 0; i < 10; i++ {
			result := db.Clauses(clauses.OrIgnore).Create(m)

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
	}

	sync := inject.MakeSync()
	return sync, inject, db
}

func TestRepo(t *testing.T) {
	m := model.Repo{
		Name: "repo2",
		URL:  "http://127.0.0.1:5000/repo2.git",
		Desc: "repo 2",
	}

	syncobj, inject, _ := setup(t, "TestRepo", &m)
	syncobj.Repo()
	assert.DirExists(t, filepath.Clean(fmt.Sprintf(`%s/%s`, inject.PathRepoDir, `127.0.0.1/repo2`)))
}

func TestFiles(t *testing.T) {
	syncobj, inject, _ := setup(t, "TestTemplates")

	contents := []byte(`
- handle: tmpl1	
  url: http://127.0.0.1:5000/tmpl1.git
  desc: Template 1
- handle: tmpl2
  url: http://127.0.0.1:5000/tmpl2.git
  desc: Template 2
`)
	err := ioutil.WriteFile(inject.PathTemplateConf, contents, 0644)
	if err != nil {
		t.Error(err)
	}

	syncobj.Files()
	assert.DirExists(t, filepath.Clean(fmt.Sprintf(`%s/%s`, inject.PathTemplateDir, `127.0.0.1/tmpl1`)))
	assert.DirExists(t, filepath.Clean(fmt.Sprintf(`%s/%s`, inject.PathTemplateDir, `127.0.0.1/tmpl2`)))
}
