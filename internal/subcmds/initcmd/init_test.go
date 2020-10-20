package initcmd

import (
	"fmt"
	fp "path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/utils"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	utils.ExitMode(utils.MPanic)
	home := fp.Join(utils.TempDir(), "init")
	set := settings.GetSettings(home)
	InitCmd([]string{"init"}, set)
	dbfile := fp.Clean(fmt.Sprintf("%s/.kick/metadata/metadata.db", home))
	assert.DirExists(t, fp.Clean(fmt.Sprintf("%s/.kick", home)))
	assert.FileExists(t, fp.Clean(fmt.Sprintf("%s/.kick/config.yml", home)))
	assert.DirExists(t, fp.Clean(fmt.Sprintf("%s/.kick/metadata", home)))
	assert.FileExists(t, dbfile)
	assert.DirExists(t, fp.Clean(fmt.Sprintf("%s/.kick/templates", home)))

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
