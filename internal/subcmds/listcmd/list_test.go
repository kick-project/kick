package listcmd_test

import (
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/kick-project/kick/internal/subcmds/listcmd"
	"github.com/kick-project/kick/internal/subcmds/updatecmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", listcmd.UsageDoc)
}

func TestList(t *testing.T) {
	exit.Mode(exit.MPanic)
	home := filepath.Join(testtools.TempDir(), "home")
	s := di.Setup(home)
	i := s.MakeSetup()
	i.Init()

	updatecmd.Update([]string{"update"}, s)

	// No entries
	args := []string{"list"}
	ret := listcmd.List(args, s)
	assert.Equal(t, 0, ret)
}

func TestListLong(t *testing.T) {
	exit.Mode(exit.MPanic)
	home := filepath.Join(testtools.TempDir(), "home")
	inject := di.Setup(home)
	i := inject.MakeSetup()
	i.Init()

	updatecmd.Update([]string{"update"}, inject)

	// No entries
	args := []string{"list", "-l"}
	ret := listcmd.List(args, inject)
	assert.Equal(t, 0, ret)
}
