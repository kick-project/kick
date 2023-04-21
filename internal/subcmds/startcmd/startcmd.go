package startcmd

import (
	"path"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `generate project scaffolding

Usage:
    kick start <handle> <project>
    kick start -s <handle>
    kick start (-l|--long)

Options:
    -h --help     print help
    -l            list templates 
    --long        list templates in long format
    -s            list template files
    <handle>      template handle
    <project>     project path
`

// OptStart start a new project from templates.
type OptStart struct {
	Start       bool   `docopt:"start"`
	Handle      string `docopt:"<handle>"`
	ProjectPath string `docopt:"<project>"`
	List        bool   `docopt:"-l"`
	ListLong    bool   `docopt:"--long"`
	Show        bool   `docopt:"-s"`
}

// Start start cli option
func Start(args []string, inject *di.DI) {
	opts := &OptStart{}
	options.Bind(UsageDoc, args, opts)
	start := inject.MakeStart()

	switch {
	case opts.List:
		start.List(false)
	case opts.Show:
		start.Show(opts.Handle, []string{}, 0)
	case opts.ListLong:
		start.List(true)
	default:
		name := path.Base(opts.ProjectPath)
		start.Start(name, opts.Handle, opts.ProjectPath)
	}
}
