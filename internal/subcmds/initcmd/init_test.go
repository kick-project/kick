package initcmd

import (
	"fmt"
	fp "path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/utils"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	home := fp.Join(utils.TempDir(), "init")
	set := settings.GetSettings(home)
	InitCmd([]string{"init"}, set)
	dbfile := fp.Clean(fmt.Sprintf("%s/.prjstart/metadata/metadata.db", home))
	assert.DirExists(t, fp.Clean(fmt.Sprintf("%s/.prjstart", home)))
	assert.FileExists(t, fp.Clean(fmt.Sprintf("%s/.prjstart/config.yml", home)))
	assert.DirExists(t, fp.Clean(fmt.Sprintf("%s/.prjstart/metadata", home)))
	assert.FileExists(t, dbfile)
	assert.DirExists(t, fp.Clean(fmt.Sprintf("%s/.prjstart/templates", home)))

	db := set.GetDB()
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
