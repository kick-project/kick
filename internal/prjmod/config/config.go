package config

import (
	"io/ioutil"

	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"gopkg.in/yaml.v2"
)

func read2bytes(f string) []byte {
	b, err := ioutil.ReadFile(f)
	errutils.Epanicf("Can not read file %s: %v", f, err) // nolint
	return b
}

type Master struct {
	Name  string      `yaml:"name"`
	Short string      `yaml:"short"`
	Desc  string      `yaml:"description"`
	Orgs  []MasterOrg `yaml:"orgs"`
}

type MasterOrg struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

func (m *Master) Load(f string) {
	b := read2bytes(f)

	err := yaml.Unmarshal(b, m)
	errutils.Epanicf("Can not unmarshal file: %s: %v", f, err)
}

type Org struct {
	Name      string   `yaml:"name"`
	Short     string   `yaml:"short"`
	Desc      string   `yaml:"desciption"`
	Templates []string `yaml:"templates"`
}

func (o *Org) Load(f string) {
	b := read2bytes(f)

	err := yaml.Unmarshal(b, o)
	errutils.Epanicf("Can not unmarshal file: %s: %v", f, err)
}

type Template struct {
	Name        string `yaml:"string"`
	Short       string `yaml:"string"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}
