package installcmd

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/apex/log"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/subcmds/initcmd"
	"github.com/crosseyed/prjstart/internal/subcmds/startcmd"
	"github.com/crosseyed/prjstart/internal/subcmds/updatecmd"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/crosseyed/prjstart/internal/utils/file"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", usageDoc)
}

func TestInstallTemplate(t *testing.T) {
	handle := "handle1"
	template := "tmpl"
	installTest(t, handle, template)
}

func TestInstallTemplateOrigin(t *testing.T) {
	handle := "handle2"
	template := "tmpl1/master1"
	installTest(t, handle, template)
}

func TestInstallTemplateURL(t *testing.T) {
	handle := "handle3"
	template := "http://localhost:5000/tmpl2.git"
	installTest(t, handle, template)
}

func installTest(t *testing.T, handle, template string) {
	id := "TestInstall"
	src := filepath.Join(utils.TempDir(), id, ".prjstart", "templates.yml.save")
	dest := filepath.Join(utils.TempDir(), id, ".prjstart", "templates.yml")
	file.Copy(src, dest)
	home := filepath.Join(utils.TempDir(), id)
	s := settings.GetSettings(home)
	s.LogLevel(log.DebugLevel)

	ec := initcmd.InitCmd([]string{"init"}, s)
	ec = updatecmd.Update([]string{"update"}, s)

	ec = Install([]string{"install", handle, template}, s)
	assert.Equal(t, 0, ec)

	td, err := ioutil.TempDir(utils.TempDir(), id+"-*")
	if err != nil {
		t.Error(err)
	}
	p := filepath.Clean(filepath.Join(td, handle))
	startcmd.Start([]string{"start", handle, p}, s)
}
