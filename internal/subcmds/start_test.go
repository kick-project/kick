package subcmds

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/config"
	"github.com/crosseyed/prjstart/internal/globals"
	"github.com/crosseyed/prjstart/internal/utils"
)

func TestStart(t *testing.T) {
	tmpdir := utils.TempDir()
	home, _ := filepath.Abs(filepath.Join(tmpdir, "home"))
	globals.Config = config.Load(home, "")
	path, _ := ioutil.TempDir(tmpdir, "start-")
	path = filepath.Join(path, "tmpl")
	args := []string{"start", "tmpl", path}
	Start(args)
}
