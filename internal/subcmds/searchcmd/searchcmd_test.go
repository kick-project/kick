package searchcmd

import (
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/services/initialize"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/utils"
)

func TestList(t *testing.T) {
	args := []string{"search", "keyword"}
	home := filepath.Join(utils.TempDir(), "home")
	s := settings.GetSettings(home)
	i := initialize.New(s.Initialize())
	i.Init()
	Search(args, s)
}
