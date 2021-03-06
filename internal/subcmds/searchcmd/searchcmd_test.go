package searchcmd_test

import (
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/kick-project/kick/internal/subcmds/searchcmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", searchcmd.UsageDoc)
}

func TestSearch(t *testing.T) {
	exit.Mode(exit.MPanic)
	args := []string{"search", "keyword"}
	home := filepath.Join(testtools.TempDir(), "home")
	inject := di.New(&di.Options{
		Home: home,
	})
	i := inject.MakeSetup()
	i.Init()
	searchcmd.Search(args, inject)
}
