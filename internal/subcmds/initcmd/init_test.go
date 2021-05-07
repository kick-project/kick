package initcmd

import (
	"fmt"
	fp "path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/utils"
	_ "github.com/mattn/go-sqlite3" // Required by 'database/sql'
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", usageDoc)
}

func TestInit(t *testing.T) {
	utils.ExitMode(utils.MPanic)
	home := fp.Join(utils.TempDir(), "init")
	inject := di.Setup(home)
	InitCmd([]string{"init"}, inject)
	dbfile := fp.Clean(fmt.Sprintf("%s/.kick/metadata/metadata.db", home))
	assert.DirExists(t, fp.Clean(fmt.Sprintf("%s/.kick", home)))
	assert.FileExists(t, fp.Clean(fmt.Sprintf("%s/.kick/config.yml", home)))
	assert.DirExists(t, fp.Clean(fmt.Sprintf("%s/.kick/metadata", home)))
	assert.FileExists(t, dbfile)
	assert.DirExists(t, fp.Clean(fmt.Sprintf("%s/.kick/templates", home)))

	db := inject.GetORM()
	row := db.Raw(`SELECT count(*) as count FROM repo WHERE url="none"`).Row()
	var count int
	err := row.Scan(&count)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, count)
}
