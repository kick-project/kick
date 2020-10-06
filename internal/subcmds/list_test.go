package subcmds

import (
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/utils"
)

func TestList(t *testing.T) {
	args := []string{"list"}
	home := filepath.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)
	ret := List(args, s)
	_ = ret
}
