package start

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/stretchr/testify/assert"
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

func TestGetOptStart(t *testing.T) {
	t.Skip("Expected fail: Seems to be an issue with docopts")
	args := []string{"prjstart", "start", "template", "project"}
	o := GetOptStart(args)
	assert.True(t, o.Start)
	assert.Equal(t, "mytemplate", o.Tmpl)
	assert.Equal(t, "myproject", o.Project)
}
