package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/crosseyed/prjstart/internal/utils/dfaults"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"gopkg.in/yaml.v2"
)

var USERCONFIG string = ".prjstart.yml"
var USERCONFIGDIR string = ".prjstart"

type Config struct {
	Home       string
	MasterURLs []string       `yaml:"masters"`
	Templates  []TemplateStub `yaml:"templates"`
}

type SortByName []TemplateStub

func (a SortByName) Len() int           { return len(a) }
func (a SortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByName) Less(i, j int) bool { return strings.Compare(a[j].Name, a[i].Name) > 0 }

type TemplateStub struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Desc string `yaml:"desc"`
}

// Load loads the main cli config config
func Load(homedir, prjstart string) *Config {
	homedir = dfaults.String(os.Getenv("HOME"), homedir)
	prjstart = dfaults.String(USERCONFIG, prjstart)
	conffile := filepath.Join(homedir, prjstart)
	if _, err := os.Stat(conffile); os.IsNotExist(err) {
		return nil
	}

	f, err := ioutil.ReadFile(conffile)
	errutils.Efatalf(err, "Can not read file %s: %v", conffile, err)
	conf := Config{
		Home: homedir,
	}

	err = yaml.Unmarshal([]byte(f), &conf)
	errutils.Efatalf(err, "Can not unmarshal file %s: %v", conffile, err)
	return &conf
}

var GLOBALCONFIG string = "prjglobal.yml"

type Global struct {
	Name        string   `yaml:"name"`
	Short       string   `yaml:"short"`
	Description string   `yaml:"description"`
	Masters     []string `yaml:"masters"`
}

func (m *Global) Load(globalconfig string) *Global {
	if globalconfig == "" {
		globalconfig = GLOBALCONFIG
	}
	if _, err := os.Stat(globalconfig); os.IsNotExist(err) {
		return nil
	}

	f, err := ioutil.ReadFile(globalconfig)
	errutils.Efatalf(err, "Can not read file %s: %v", globalconfig, err)

	err = yaml.Unmarshal([]byte(f), m)
	errutils.Efatalf(err, "Can not unmarshal file %s: %v", globalconfig, err)
	return m
}

var MASTERCONFIG string = "prjmaster.yml"

type Master struct {
	Name        string   `yaml:"name"`
	Short       string   `yaml:"short"`
	Description string   `yaml:"description"`
	Templates   []string `yaml:"templates"`
}

func (o *Master) Load(masterconfig string) *Master {
	if masterconfig == "" {
		masterconfig = MASTERCONFIG
	}
	if _, err := os.Stat(masterconfig); os.IsNotExist(err) {
		return nil
	}

	f, err := ioutil.ReadFile(masterconfig)
	errutils.Efatalf(err, "Can not read file %s: %v", masterconfig, err)

	err = yaml.Unmarshal([]byte(f), o)
	errutils.Efatalf(err, "Can not unmarshal file %s: %v", masterconfig, err)
	return o
}
