package installcmd

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/apex/log"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/subcmds/initcmd"
	"github.com/kick-project/kick/internal/subcmds/startcmd"
	"github.com/kick-project/kick/internal/subcmds/updatecmd"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/file"
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

func TestInstallPath(t *testing.T) {
	handle := "handle4"
	template := filepath.Clean(utils.TempDir() + "/TestInstall/kicks/go")
	installTest(t, handle, template)
}

func installTest(t *testing.T, handle, template string) {
	utils.ExitMode(utils.MPanic)
	id := "TestInstall"
	src := filepath.Join(utils.TempDir(), id, ".kick", "templates.yml.save")
	dest := filepath.Join(utils.TempDir(), id, ".kick", "templates.yml")
	_, err := file.Copy(src, dest)
	if err != nil {
		t.Error(err)
	}

	home := filepath.Join(utils.TempDir(), id)
	s := settings.GetSettings(home)
	s.LogLevel(log.DebugLevel)

	ec := initcmd.InitCmd([]string{"init"}, s)
	assert.Equal(t, 0, ec)

	ec = updatecmd.Update([]string{"update"}, s)
	assert.Equal(t, 0, ec)

	ec = Install([]string{"install", handle, template}, s)
	assert.Equal(t, 0, ec)

	td, err := ioutil.TempDir(utils.TempDir(), id+"-*")
	if err != nil {
		t.Error(err)
	}
	p := filepath.Clean(filepath.Join(td, handle))
	startcmd.Start([]string{"start", handle, p}, s)
}
