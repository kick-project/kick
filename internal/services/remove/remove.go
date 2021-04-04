package remove

import (
	"fmt"
	"io"

	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/utils/errutils"
)

// Remove remove installed templates
type Remove struct {
	Conf             *config.File
	PathTemplateConf string
	PathUserConf     string
	Stderr           io.Writer
	Stdout           io.Writer
}

// Remove removes a handle from installed templates
func (r *Remove) Remove(handle string) int {
	item := -1
	templates := r.Conf.Templates
	l := len(templates)
	for i, t := range templates {
		if t.Handle != handle {
			continue
		}
		item = i
		break
	}

	if item == -1 {
		fmt.Fprintf(r.Stderr, "can not uninstall handle %s. handle not installed\n", handle)
		return 255
	}

	switch item {
	case 0: // beginning
		if l > 1 {
			templates = templates[item+1:]
		} else {
			templates = []config.Template{}
		}
	case l - 1: // end
		templates = templates[0:item]
	default: // middle
		head := templates[0:item]
		tail := templates[item+1:]

		templates = append(head, tail...)
	}
	r.Conf.Templates = templates
	err := r.Conf.SaveTemplates()
	errutils.Epanic(err)

	return 0
}
