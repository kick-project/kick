package config

import (
	// "github.com/stretchr/testify/assert"
	"path"
	"testing"

	"github.com/crosseyed/prjstart/internal/utils"
)

func TestMasterLoad(t *testing.T) {
	fpath := path.Join(utils.FixtureDir(), "prjmod", "master.yaml")
	d := Master{}
	d.Load(fpath)
}

func TestOrgLoad(t *testing.T) {
	fpath := path.Join(utils.FixtureDir(), "prjmod", "org.yaml")
	d := Master{}
	d.Load(fpath)
}

func TestTemplateLoad(t *testing.T) {
	fpath := path.Join(utils.FixtureDir(), "prjmod", "template.yaml")
	d := Master{}
	d.Load(fpath)
}
