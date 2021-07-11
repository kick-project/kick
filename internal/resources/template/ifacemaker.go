// DO NOT EDIT: Generated using "make interfaces"

package template

import (
	"github.com/kick-project/kick/internal/resources/template/variables"
)

// TemplateIface ...
type TemplateIface interface {
	// SetRender set rendering engine
	SetRender(renderer string)
	SetVars(vars *variables.Variables)
	// SetSrcDest sets the source template and destination path where the project structure
	// will reside.
	SetSrcDest(src, dest string)
	// SetSrc sets the source template "name". "name" is defined
	// in *config.Config.TemplateURLs. *config.Config is provided as an Option to New.
	SetSrc(name string)
	// SetDest sets the destination path
	SetDest(dest string)
	// Run generates the target directory structure
	Run() int
}
