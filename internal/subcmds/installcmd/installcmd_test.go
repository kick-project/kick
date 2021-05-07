package installcmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/apex/log"
	"github.com/kick-project/kick/internal/di"
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
	installTest(t, "TestInstallTemplate", handle, template)
}

func TestInstallTemplateOrigin(t *testing.T) {
	handle := "handle2"
	template := "tmpl1/repo1"
	installTest(t, "TestInstallTemplateOrigin", handle, template)
}

func TestInstallTemplateURL(t *testing.T) {
	handle := "handle3"
	template := "http://localhost:5000/tmpl2.git"
	installTest(t, "TestInstallTemplateURL", handle, template)
}

func TestInstallPath(t *testing.T) {
	handle := "handle4"
	template := filepath.Clean(utils.TempDir() + "/installcmd/kicks/go")
	installTest(t, "TestInstallPath", handle, template)
}

func installTest(t *testing.T, id, handle, template string) {
	utils.ExitMode(utils.MPanic)
	// Home Directory
	home := filepath.Join(utils.TempDir(), id)

	// Make kick config dir
	kickDir := filepath.Join(home, ".kick")
	err := os.MkdirAll(kickDir, 0755)
	if err != nil {
		t.Errorf("Can not create directory \"%s\": %v", kickDir, err)
		return
	}

	// Copy template
	srcTemplate := filepath.Join(utils.FixtureDir(), "installcmd", ".kick", "templates.yml.save")
	destTemplate := filepath.Join(kickDir, "templates.yml")
	_, err = file.Copy(srcTemplate, destTemplate)
	if err != nil {
		t.Error(err)
	}

	// Copy config
	srcConfig := filepath.Join(utils.FixtureDir(), "installcmd", ".kick", "config.yml")
	destConfig := filepath.Join(kickDir, "config.yml")
	_, err = file.Copy(srcConfig, destConfig)
	if err != nil {
		t.Error(err)
	}

	inject := di.Setup(home)
	inject.LogLevel(log.DebugLevel)

	ec := initcmd.InitCmd([]string{"init"}, inject)
	assert.Equal(t, 0, ec)

	ec = updatecmd.Update([]string{"update"}, inject)
	assert.Equal(t, 0, ec)

	ec = Install([]string{"install", handle, template}, inject)
	assert.Equal(t, 0, ec)

	td, err := ioutil.TempDir(utils.TempDir(), id+"-*")
	if err != nil {
		t.Error(err)
	}
	p := filepath.Clean(filepath.Join(td, handle))
	ec = startcmd.Start([]string{"start", handle, p}, inject)
	assert.Equal(t, 0, ec)
}
