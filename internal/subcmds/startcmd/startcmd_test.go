package startcmd_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/kick-project/kick/internal/subcmds/startcmd"
	"github.com/stretchr/testify/assert"
)

func TestUsageDoc(t *testing.T) {
	assert.NotRegexp(t, "\t", startcmd.UsageDoc)
}

func TestStart(t *testing.T) {
	exit.Mode(exit.MPanic)
	tmpdir := testtools.TempDir()
	home, _ := filepath.Abs(filepath.Join(tmpdir, "home"))
	path, _ := ioutil.TempDir(tmpdir, "start-")
	project := "tmpl"
	tmpl := "tmpl"
	path = filepath.Join(path, project)
	args := []string{"start", tmpl, path}
	inject := di.New(&di.Options{
		Home: home,
	})
	startcmd.Start(args, inject)

	type interpolated struct {
		Project string `yaml:"project"`
		Home    string `yaml:"home"`
	}
	d := &interpolated{}
	y := filepath.Join(path, "template-interpolate.yml")
	err := marshal.FromFile(d, y)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, d.Project, project)
	assert.NotEmpty(t, os.Getenv("HOME"))
	assert.Equal(t, d.Home, os.Getenv("HOME"))
}
