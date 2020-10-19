package searchcmd

import (
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/iinitialize"
	"github.com/kick-project/kick/internal/utils"
	"github.com/jinzhu/copier"
)

func TestSearch(t *testing.T) {
	args := []string{"search", "keyword"}
	home := filepath.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)
	i := &initialize.Initialize{}
	copier.Copy(i, iinitialize.Inject(s))
	i.Init()
	Search(args, s)
}
