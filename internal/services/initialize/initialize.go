// Package initialize initializes a template or a repository
package initialize

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kick-project/kick/internal/resources/config/configtemplate"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/serialize"
)

// Init create repositories and templates
//
//go:generate ifacemaker -f initialize.go -s Init -p initialize -i InitIface -o initialize_interfaces.go -c "AUTO GENERATED. DO NOT EDIT."
type Init struct {
	ErrHandler errs.HandlerIface  `validate:"required"`
	Log        logger.OutputIface `validate:"required"`
}

// CreateRepo create repository
func (i *Init) CreateRepo(name, path string) int {
	var (
		wd  string
		err error
	)
	if path != "" {
		err = os.Mkdir(name, 0755)
		if i.ErrHandler.LogF(`can not create repo "%s": %v`, name, err) {
			return 255
		}
		wd = name
	} else {
		wd = "."
	}

	repo := &serialize.RepoMain{
		Name: name,
		Desc: fmt.Sprintf(`Repository %s`, name),
		TemplateURLs: []string{
			`http://example.template.host/template.git`,
		},
	}
	err = marshal.ToFile(repo, filepath.Join(wd, "repo.yml"))
	if i.ErrHandler.LogF(`can not create repo "%s": %v`, name, err) {
		return 255
	}

	if path == "" {
		i.Log.Printf(`generated %s`, `repo.yml`)
	} else {
		i.Log.Printf(`generated %s`, filepath.Join(path, `repo.yml`))
	}
	return 0
}

// CreateTemplate create template
func (i *Init) CreateTemplate(name, path string) int {
	var (
		wd  string
		err error
	)
	if path != "" {
		err := os.Mkdir(name, 0755)
		if i.ErrHandler.LogF(`can not create repo "%s": %v`, name, err) {
			return 255
		}
		wd = name
	} else {
		wd, err = os.Getwd()
		i.ErrHandler.FatalF(`can not find current directory: %v`, err)
	}

	tmpl := &configtemplate.TemplateMain{
		Name: name,
		Desc: fmt.Sprintf(`Template %s`, name),
	}
	err = marshal.ToFile(tmpl, filepath.Join(wd, ".kick.yml"))
	if i.ErrHandler.LogF(`can not create template "%s": %v`, name, err) {
		return 255
	}

	if path == "" {
		i.Log.Printf(`created %s`, `.kick.yml`)
	} else {
		i.Log.Printf(`created %s`, filepath.Join(path, `.kick.yml`))
	}
	return 0
}
