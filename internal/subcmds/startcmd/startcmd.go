package startcmd

import (
	"fmt"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/resources/template"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/icheck"
	"github.com/kick-project/kick/internal/settings/isync"
	"github.com/kick-project/kick/internal/settings/itemplate"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `Generate project scaffolding

Usage:
    kick start <handle> <project>

Options:
    -h --help     Print help.
    <handle>      Template handle.
    <project>     Project path.
`

// OptStart start a new project from templates.
type OptStart struct {
	Start       bool   `docopt:"start"`
	Template    string `docopt:"<handle>"`
	ProjectPath string `docopt:"<project>"`
}

// Start start cli option
func Start(args []string, s *settings.Settings) int {
	opts := &OptStart{}
	options.Bind(usageDoc, args, opts)

	chk := &check.Check{}
	err := copier.Copy(chk, icheck.Inject(s))
	errutils.Epanic(err)

	if err = chk.Init(); err != nil {
		fmt.Fprintf(s.Stderr, "%s\n", err.Error())
		utils.Exit(255)
	}

	// Sync DB table "installed" with configuration file
	synchro := &sync.Sync{}
	err = copier.Copy(synchro, isync.Inject(s))
	errutils.Epanic(err)
	synchro.Templates()

	// Set project name
	s.ProjectName = filepath.Base(opts.ProjectPath)
	t := &template.Template{}
	err = copier.Copy(t, itemplate.Inject(s))
	errutils.Epanic(err)
	t.SetSrcDest(opts.Template, opts.ProjectPath)
	ret := t.Run()
	return ret
}
