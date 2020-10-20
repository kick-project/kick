package updatecmd

import (
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/subcmds/initcmd"
	"github.com/kick-project/kick/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", usageDoc)
}

func TestUpdate(t *testing.T) {
	utils.ExitMode(utils.MPanic)

	home := filepath.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)

	initcmd.InitCmd([]string{"init"}, s)

	dbConn := s.GetDB()
	_, err := dbConn.Exec(`DELETE FROM templates`)
	if err != nil {
		t.Error(err)
	}

	Update([]string{"update"}, s)

	var count int
	row := dbConn.QueryRow(`SELECT count(*) AS count FROM templates`)
	err = row.Scan(&count)
	if err != nil {
		t.Error(err)
	}

	assert.GreaterOrEqual(t, count, 1)
}
