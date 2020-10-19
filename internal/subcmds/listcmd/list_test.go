package listcmd

import (
	"path/filepath"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/iinitialize"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
)

func TestList(t *testing.T) {
	args := []string{"list"}
	home := filepath.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)
	i := &initialize.Initialize{}
	err := copier.Copy(i, iinitialize.Inject(s))
	errutils.Epanic(err)
	i.Init()
	ret := List(args, s)
	_ = ret
}
