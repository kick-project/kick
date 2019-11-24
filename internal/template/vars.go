package template

import (
	"os"
	"strings"
)

// TmplVars consists of the variables that are to be passed to the template
type TmplVars struct {
	Env     map[string]string
	Project map[string]string
}

// NewTmplVars sets up Environment and project variables to be passed through to the text template
func NewTmplVars(prj map[string]string) *TmplVars {
	tv := TmplVars{}
	tv.genVarsEnv()
	tv.Project = prj
	return &tv
}

func (v *TmplVars) genVarsEnv() map[string]string {
	envMap := make(map[string]string)

	for _, v := range os.Environ() {
		split_v := strings.Split(v, "=")
		envMap[split_v[0]] = split_v[1]
	}
	v.Env = envMap

	return envMap
}
