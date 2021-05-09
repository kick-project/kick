package startcmd

import (
	"fmt"
	"path/filepath"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/utils/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `generate project scaffolding

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
	options.Bind(UsageDoc, args, opts)

	chk := inject.MakeCheck()

	if err := chk.Init(); err != nil {
		fmt.Fprintf(inject.Stderr, "%s\n", err.Error())
		exit.Exit(255)
	}

	// Sync DB table "installed" with configuration file
	synchro := inject.MakeSync()
	synchro.Files()

	// Set project name
	inject.ProjectName = filepath.Base(opts.ProjectPath)
	t := inject.MakeTemplate()
	t.SetSrcDest(opts.Template, opts.ProjectPath)
	ret := t.Run()
	return ret
}
