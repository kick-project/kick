package startcmd

import (
	"path/filepath"

	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/resources/template"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/isync"
	"github.com/kick-project/kick/internal/settings/itemplate"
	"github.com/kick-project/kick/internal/utils/options"
	"github.com/jinzhu/copier"
)

var usageDoc = `Generate project scaffolding

Usage:
    prjstart start <handle> <project>

Options:
    -h --help     Print help.
    <handle>      Template handle.
    <project>     Project path.
`

type OptStart struct {
	Start       bool   `docopt:"start"`
	Template    string `docopt:"<handle>"`
	ProjectPath string `docopt:"<project>"`
}

// Start start cli option
func Start(args []string, s *settings.Settings) int {
	opts := &OptStart{}
	options.Bind(usageDoc, args, opts)

	// Sync DB table "installed" with configuration file
	synchro := &sync.Sync{}
	copier.Copy(synchro, isync.Inject(s))
	synchro.Templates()

	// Set project name
	s.ProjectName = filepath.Base(opts.ProjectPath)
	t := &template.Template{}
	copier.Copy(t, itemplate.Inject(s))
	t.SetSrcDest(opts.Template, opts.ProjectPath)
	ret := t.Run()
	return ret
}
