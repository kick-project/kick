package startcmd

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", usageDoc)
}

func TestStart(t *testing.T) {
	exit.Mode(exit.MPanic)
	tmpdir := utils.TempDir()
	home, _ := filepath.Abs(filepath.Join(tmpdir, "home"))
	path, _ := ioutil.TempDir(tmpdir, "start-")
	path = filepath.Join(path, "tmpl")
	args := []string{"start", "tmpl", path}
	inject := di.Setup(home)
	Start(args, inject)
	assert.DirExists(t, path)
}
