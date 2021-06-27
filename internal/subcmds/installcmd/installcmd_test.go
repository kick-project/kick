package installcmd_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/file"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/kick-project/kick/internal/subcmds/installcmd"
	"github.com/kick-project/kick/internal/subcmds/setupcmd"
	"github.com/kick-project/kick/internal/subcmds/startcmd"
	"github.com/kick-project/kick/internal/subcmds/updatecmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", installcmd.UsageDoc)
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

func TestInstallSSHPublic(t *testing.T) {
	if os.Getenv("KICK_TEST_SSH") != "true" {
		t.Skip(`SSH git repo tests are disabled`)
	}
	handle := "public1"
	template := "git@github.com:kick-fixtures/template-public.git"
	installTest(t, "TestInstallSSHPublic", handle, template)
}

func TestInstallSSHPrivate(t *testing.T) {
	if os.Getenv("KICK_TEST_SSH") != "true" {
		t.Skip(`SSH git repo tests are disabled`)
	}
	handle := "private1"
	template := "git@github.com:kick-fixtures/template-private.git"
	installTest(t, "TestInstallSSHPrivate", handle, template)
}

func TestInstallSSHNoRepo(t *testing.T) {
	handle := "norepo1"
	template := "git@github.com:kick-fixtures/template-norepo.git"
	id := "TestInstallSSHNoRepo"

	// Home Directory
	home := filepath.Join(testtools.TempDir(), id)

	// Make kick config dir
	kickDir := filepath.Join(home, ".kick")
	err := os.MkdirAll(kickDir, 0755)
	if err != nil {
		t.Errorf("Can not create directory \"%s\": %v", kickDir, err)
		return
	}

	inject := di.New(&di.Options{Home: home})
	inject.LogLevel(logger.DebugLevel)
	inject.ExitMode = exit.MPanic

	ec := setupcmd.SetupCmd([]string{"setup"}, inject)
	assert.Equal(t, 0, ec)

	ec = updatecmd.Update([]string{"update"}, inject)
	assert.Equal(t, 0, ec)

	ec = installcmd.Install([]string{"install", handle, template}, inject)
	assert.NotEqual(t, 0, ec)
}

func TestInstallPath(t *testing.T) {
	handle := "handle4"
	template := filepath.Clean(testtools.TempDir() + "/installcmd/kicks/go")
	installTest(t, "TestInstallPath", handle, template)
}

func installTest(t *testing.T, id, handle, template string) {
	exit.Mode(exit.MPanic)
	// Home Directory
	home := filepath.Join(testtools.TempDir(), id)

	// Make kick config dir
	kickDir := filepath.Join(home, ".kick")
	err := os.MkdirAll(kickDir, 0755)
	if err != nil {
		t.Errorf("Can not create directory \"%s\": %v", kickDir, err)
		return
	}

	// Copy template
	srcTemplate := filepath.Join(testtools.FixtureDir(), "installcmd", ".kick", "templates.yml.save")
	destTemplate := filepath.Join(kickDir, "templates.yml")
	_, err = file.Copy(srcTemplate, destTemplate)
	if err != nil {
		t.Error(err)
	}

	// Copy config
	srcConfig := filepath.Join(testtools.FixtureDir(), "installcmd", ".kick", "config.yml")
	destConfig := filepath.Join(kickDir, "config.yml")
	_, err = file.Copy(srcConfig, destConfig)
	if err != nil {
		t.Error(err)
	}

	inject := di.New(&di.Options{Home: home})
	inject.LogLevel(logger.DebugLevel)

	ec := setupcmd.SetupCmd([]string{"setup"}, inject)
	assert.Equal(t, 0, ec)

	ec = updatecmd.Update([]string{"update"}, inject)
	assert.Equal(t, 0, ec)

	ec = installcmd.Install([]string{"install", handle, template}, inject)
	assert.Equal(t, 0, ec)

	td, err := ioutil.TempDir(testtools.TempDir(), id+"-*")
	if err != nil {
		t.Error(err)
	}
	p := filepath.Clean(filepath.Join(td, handle))
	ec = startcmd.Start([]string{"start", handle, p}, inject)
	assert.Equal(t, 0, ec)
}
