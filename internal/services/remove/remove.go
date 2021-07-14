package remove

import (
	"io"

	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/logger"
)

// Remove remove installed templates
//go:generate ifacemaker -f remove.go -s Remove -p remove -i RemoveIface -o remote_interfaces.go -c "AUTO GENERATED. DO NOT EDIT."
type Remove struct {
	Conf             *config.File       `validate:"required"`
	Err              *errs.Handler      `validate:"required"`
	Log              logger.OutputIface `validate:"required"`
	PathTemplateConf string             `validate:"required"`
	PathUserConf     string             `validate:"required"`
	Stderr           io.Writer          `validate:"required"`
	Stdout           io.Writer          `validate:"required"`
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
		r.Log.Printf("can not uninstall handle %s. handle not installed\n", handle)
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
	r.Err.Panic(err)
	r.Log.Printf(`removed handle:%s`, handle)
	return 0
}
