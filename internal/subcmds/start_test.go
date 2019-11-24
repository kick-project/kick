package subcmds

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestStart(t *testing.T) {
	tmpdir, _ := filepath.Abs("../../tmp")
	path, _ := ioutil.TempDir(tmpdir, "start-")
	path = filepath.Join(path, "tmpl")
	args := []string{"start", "tmpl", path}
	Start(args)
}
