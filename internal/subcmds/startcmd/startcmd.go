package startcmd

import (
	"fmt"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/icheck"
	"github.com/kick-project/kick/internal/di/isync"
	"github.com/kick-project/kick/internal/di/itemplate"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/resources/template"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `generate project scaffolding

Usage:
    kick start <handle> <project>

Options:
    -h --help     print help
    <handle>      template handle
    <project>     project path
`

// OptStart start a new project from templates.
type OptStart struct {
	Start       bool   `docopt:"start"`
	Template    string `docopt:"<handle>"`
	ProjectPath string `docopt:"<project>"`
}

// Start start cli option
func Start(args []string, inject *di.DI) int {
	opts := &OptStart{}
	options.Bind(usageDoc, args, opts)

	chk := &check.Check{}
	err := copier.Copy(chk, icheck.Inject(inject))
	errutils.Epanic(err)

	if err = chk.Init(); err != nil {
		fmt.Fprintf(inject.Stderr, "%s\n", err.Error())
		exit.Exit(255)
	}

	// Sync DB table "installed" with configuration file
	synchro := &sync.Sync{}
	err = copier.Copy(synchro, isync.Inject(inject))
	errutils.Epanic(err)
	synchro.Files()

	// Set project name
	inject.ProjectName = filepath.Base(opts.ProjectPath)
	t := &template.Template{}
	err = copier.Copy(t, itemplate.Inject(inject))
	errutils.Epanic(err)
	t.SetSrcDest(opts.Template, opts.ProjectPath)
	ret := t.Run()
	return ret
}
