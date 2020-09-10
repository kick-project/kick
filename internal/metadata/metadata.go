package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/crosseyed/prjstart/internal/file"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"gopkg.in/yaml.v2"
)

type TemplateStub struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Desc string `yaml:"desc"`
}

var GLOBALMETA string = "metadata.global.yml"

// Global is the global metadata
type Global struct {
	Name        string   `yaml:"name"`
	URL         string   `yaml:"url"`
	Description string   `yaml:"description"`
	Masters     []string `yaml:"masters"`
}

// Build builds the metadata from configuration file settings
func (g *Global) Build(configfile string) {
}

// Load loads from a yamlfile
func (g *Global) Load(yamlfile string) {
	if yamlfile == "" {
		yamlfile = GLOBALMETA
	}
	if _, err := os.Stat(yamlfile); os.IsNotExist(err) {
		return
	}

	f, err := ioutil.ReadFile(yamlfile)
	errutils.Efatalf("Can not read file %s: %v", yamlfile, err)

	err = yaml.Unmarshal([]byte(f), g)
	errutils.Efatalf("Can not unmarshal file %s: %v", yamlfile, err)
	return
}

// Save saves to yamlfile.
func (g *Global) Save(yamlfile string) {
	if yamlfile == "" {
		yamlfile = GLOBALMETA
	}
	d, err := filepath.Abs(filepath.Dir(yamlfile))
	errutils.Epanicf("Can not get absolute path: %w", err)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		errutils.Efatalf("Parent directory of %s does not exists: %w", yamlfile, err)
	}

	f := file.NewAtomicWrite(yamlfile)
	defer f.Close()
	out, err := yaml.Marshal(g)
	errutils.Efatalf("Can not marshal: %w", err)
	f.Write(out)
	errutils.Efatalf("Can not write to file: %w", err)
}

var MASTERMETA string = "metadata.master.yml"

type Master struct {
	Name        string   `yaml:"name"`
	URL         string   `yaml:"url"`
	Description string   `yaml:"description"`
	Templates   []string `yaml:"templates"`
}

// Build builds the metadata from configuration file settings
func (m *Master) Build(configfile string) {
}

// Load loads from a yaml file
func (m *Master) Load(yamlfile string) {
	if yamlfile == "" {
		yamlfile = MASTERMETA
	}
	if _, err := os.Stat(yamlfile); os.IsNotExist(err) {
		return
	}

	f, err := ioutil.ReadFile(yamlfile)
	errutils.Efatalf("Can not read file %s: %v", yamlfile, err)

	err = yaml.Unmarshal([]byte(f), m)
	errutils.Efatalf("Can not unmarshal file %s: %v", yamlfile, err)
}

// Save saves to yamlfile.
func (m *Master) Save(yamlfile string) {
	if yamlfile == "" {
		yamlfile = GLOBALMETA
	}
	d, err := filepath.Abs(filepath.Dir(yamlfile))
	errutils.Epanicf("Can not get absolute path: %w", err)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		errutils.Efatalf("Parent directory of %s does not exists: %w", yamlfile, err)
	}

	f := file.NewAtomicWrite(yamlfile)
	defer f.Close()
	out, err := yaml.Marshal(m)
	errutils.Efatalf("Can not marshal: %w", err)
	f.Write(out)
	errutils.Efatalf("Can not write to file: %w", err)
}
