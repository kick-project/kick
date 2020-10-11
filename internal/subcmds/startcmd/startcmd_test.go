package startcmd

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/utils"
)

func TestStart(t *testing.T) {
	tmpdir := utils.TempDir()
	home, _ := filepath.Abs(filepath.Join(tmpdir, "home"))
	path, _ := ioutil.TempDir(tmpdir, "start-")
	path = filepath.Join(path, "tmpl")
	args := []string{"start", "tmpl", path}
	s := settings.GetSettings(home)
	Start(args, s)
}
