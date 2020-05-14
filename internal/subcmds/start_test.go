package subcmds

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/config"
	"github.com/crosseyed/prjstart/internal/globals"
)

func TestStart(t *testing.T) {
	tmpdir, _ := filepath.Abs("../../tmp")
	home, _ := filepath.Abs("../../tmp/home")
	globals.Config = config.Load(home, "")
	path, _ := ioutil.TempDir(tmpdir, "start-")
	path = filepath.Join(path, "tmpl")
	args := []string{"start", "tmpl", path}
	Start(args)
}
