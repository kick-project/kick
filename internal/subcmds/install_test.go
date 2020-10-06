package subcmds

import (
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/utils"
)

func TestInstall(t *testing.T) {
	args := []string{"install", "set1", "tmpl1"}
	home := filepath.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)
	Install(args, s)
}
