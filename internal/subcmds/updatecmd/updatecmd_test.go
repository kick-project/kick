package updatecmd_test

import (
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/kick-project/kick/internal/subcmds/setupcmd"
	"github.com/kick-project/kick/internal/subcmds/updatecmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", updatecmd.UsageDoc)
}

func TestUpdate(t *testing.T) {
	exit.Mode(exit.MPanic)

	home := filepath.Join(testtools.TempDir(), "home")
	inject := di.Setup(home)

	setupcmd.SetupCmd([]string{"setup"}, inject)

	dbConn := inject.MakeORM()
	result := dbConn.Raw(`DELETE FROM template`)
	if result.Error != nil {
		t.Error(result.Error)
	}

	updatecmd.Update([]string{"update"}, inject)

	var count int
	row := dbConn.Raw(`SELECT count(*) AS count FROM template`).Row()
	err := row.Scan(&count)
	if err != nil {
		t.Error(err)
	}

	assert.GreaterOrEqual(t, count, 1)
}
