package internal

import (
	"github.com/crosseyed/prjstart/internal/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ConfigStruct struct {
	Home      string
	SetURLs   []string       `yaml:"sets"`
	Templates []TemplateStub `yaml:"templates"`
}

type SortByName []TemplateStub

func (a SortByName) Len() int           { return len(a) }
func (a SortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByName) Less(i, j int) bool { return strings.Compare(a[j].Name, a[i].Name) > 0 }

type SetList struct {
	Set         string         `yaml:"set"`
	Description string         `yaml:"desc"`
	Templates   []TemplateStub `yaml:"templates"`
}

type TemplateStub struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Desc string `yaml:"desc"`
	Set  string `yaml:"set"`
}

func LoadConfig(homedir, prjstart string) *ConfigStruct {
	if homedir == "" {
		homedir = os.Getenv("HOME")
	}
	if prjstart == "" {
		prjstart = ".prjstart.yml"
	}
	conffile := filepath.Join(homedir, prjstart)
	if _, err := os.Stat(conffile); os.IsNotExist(err) {
		return nil
	}

	f, err := ioutil.ReadFile(conffile)
	utils.ChkErr(err, utils.Efatalf, "Can not read file %s: %v", conffile, err)
	conf := ConfigStruct{
		Home: homedir,
	}

	err = yaml.Unmarshal([]byte(f), &conf)
	utils.ChkErr(err, utils.Efatalf, "Can not unmarshal file %s: %v", conffile, err)
	return &conf
}

func LoadSet(path, conffile string) *SetList {
	if conffile == "" {
		conffile = ".prjstart.yml"
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	confpath := filepath.Join(path, conffile)
	f, err := ioutil.ReadFile(confpath)
	utils.ChkErr(err, utils.Efatalf, "Can not read file %s: %v", path, err)
	conf := SetList{}

	err = yaml.Unmarshal([]byte(f), &conf)
	utils.ChkErr(err, utils.Efatalf, "Can not unmarshal file %s: %v", path, err)
	return &conf
}
