package listcmd

import (
	"path/filepath"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/iinitialize"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/subcmds/updatecmd"
	"github.com/kick-project/kick/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", usageDoc)
}

func TestList(t *testing.T) {
	exit.Mode(exit.MPanic)
	home := filepath.Join(utils.TempDir(), "home")
	s := di.Setup(home)
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
	exit.Mode(exit.MPanic)
	home := filepath.Join(utils.TempDir(), "home")
	inject := di.Setup(home)
	i := &initialize.Initialize{}
	err := copier.Copy(i, iinitialize.Inject(inject))
	if err != nil {
		t.Error(err)
	}
	i.Init()

	updatecmd.Update([]string{"update"}, inject)

	// No entries
	args := []string{"list", "-l"}
	ret := List(args, inject)
	assert.Equal(t, 0, ret)
}
