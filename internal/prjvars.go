package internal

import (
	"path/filepath"

	"github.com/crosseyed/prjstart/internal/services/template"
)

type PrjVars [][]string

func SetVars(opts *OptStart) PrjVars {
	data := [][]string{
		{"NAME", filepath.Base(opts.Project), "Project Name derived from target directory"},
	}
	return data
}

func (prjvars PrjVars) GetVars() *template.Variables {
	vars := template.NewTmplVars()
	for _, v := range prjvars {
		vars.SetProjectVar(v[0], v[1])
	}
	return vars
}

func (prjvars PrjVars) GetDescriptions() [][]string {
	data := [][]string{}
	for _, v := range prjvars {
		data = append(data, []string{v[0], v[2]})
	}
	return data
}
