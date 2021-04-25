package searchcmd

import (
	"path/filepath"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/iinitialize"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", usageDoc)
}

func TestSearch(t *testing.T) {
	utils.ExitMode(utils.MPanic)
	args := []string{"search", "keyword"}
	home := filepath.Join(utils.TempDir(), "home")
	inject := di.Setup(home)
	i := &initialize.Initialize{}
	err := copier.Copy(i, iinitialize.Inject(inject))
	errutils.Epanic(err)
	i.Init()
	Search(args, inject)
}
