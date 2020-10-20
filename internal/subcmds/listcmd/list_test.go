package listcmd

import (
	"path/filepath"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/iinitialize"
	"github.com/kick-project/kick/internal/subcmds/updatecmd"
	"github.com/kick-project/kick/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", usageDoc)
}

func TestList(t *testing.T) {
	utils.ExitMode(utils.MPanic)
	home := filepath.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)
	i := &initialize.Initialize{}
	err := copier.Copy(i, iinitialize.Inject(s))
	if err != nil {
		t.Error(err)
	}
	i.Init()

	updatecmd.Update([]string{"update"}, s)

	// No entries
	args := []string{"list"}
	ret := List(args, s)
	assert.Equal(t, 0, ret)
}

func TestListLong(t *testing.T) {
	utils.ExitMode(utils.MPanic)
	home := filepath.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)
	i := &initialize.Initialize{}
	err := copier.Copy(i, iinitialize.Inject(s))
	if err != nil {
		t.Error(err)
	}
	i.Init()

	updatecmd.Update([]string{"update"}, s)

	// No entries
	args := []string{"list", "-l"}
	ret := List(args, s)
	assert.Equal(t, 0, ret)
}
