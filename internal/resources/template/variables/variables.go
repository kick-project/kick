package variables

import (
	"os"
	"strings"
)

// Variables consists of the variables that are to be passed to the template
type Variables struct {
	Env     map[string]string // Environment variables
	Project map[string]string // Project variables
}

// New sets up environment and project variables to be passed through to the text template
func New() *Variables {
	tv := Variables{}
	tv.genVarsEnv()
	tv.Project = map[string]string{}
	return &tv
}

// ProjectVariable sets a project variable
func (v *Variables) ProjectVariable(name, value string) {
	v.Project[name] = value
	os.Setenv("PROJECT_"+name, value)
}

func (v *Variables) genVarsEnv() map[string]string {
	envMap := make(map[string]string)

	for _, v := range os.Environ() {
		splitVars := strings.Split(v, "=")
		envMap[splitVars[0]] = splitVars[1]
	}
	v.Env = envMap

	return envMap
}
