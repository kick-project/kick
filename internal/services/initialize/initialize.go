// Package initialize initializes a template or a repository
package initialize

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/serialize"
)

// Init create repositories and templates
type Init struct {
	ErrHandler errs.HandlerIface `validate:"required"`
	Stdout     io.Writer         `validate:"required"`
	Stderr     io.Writer         `validate:"required"`
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

	tmpl := &serialize.TemplateMain{
		Name: name,
		Desc: fmt.Sprintf(`Template %s`, name),
	}
	err = marshal.ToFile(tmpl, filepath.Join(wd, ".kick.yml"))
	if i.ErrHandler.LogF(`can not create template "%s": %v`, name, err) {
		return 255
	}
	return 0
}
