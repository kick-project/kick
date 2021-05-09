package searchcmd_test

import (
	"path/filepath"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/iinitialize"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/subcmds/searchcmd"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", searchcmd.UsageDoc)
}

func TestSearch(t *testing.T) {
	exit.Mode(exit.MPanic)
	args := []string{"search", "keyword"}
	home := filepath.Join(utils.TempDir(), "home")
	inject := di.Setup(home)
	i := &initialize.Initialize{}
	err := copier.Copy(i, iinitialize.Inject(inject))
	errutils.Epanic(err)
	i.Init()
	searchcmd.Search(args, inject)
}
