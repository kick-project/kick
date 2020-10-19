package startcmd

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/utils"
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
