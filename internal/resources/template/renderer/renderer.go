package renderer

import (
	"regexp"

	"github.com/kick-project/kick/internal/resources/template/variables"
)

// Renderer is an interface to render templates.
type Renderer interface {
	File2File(src, dst string, vars *variables.Variables, nounset, noempty bool) error
	Text2File(text, dst string, vars *variables.Variables, nounset, noempty bool) error
	Text2String(text string, vars *variables.Variables, nounset, noempty bool) (string, error)
	RenderDirRegexp() *regexp.Regexp
}
