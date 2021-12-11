package initcmd

import (
	"errors"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `Create a repo or template

Usage:
    kick init repo <name> [<path>]
    kick init template <name> [<path>]

Options:
    -h --help    print help
    repo         create repository
    template     create a template
    <name>       template or repo name
    <path>       directory path. if not set creates files in working directory
`

// OptInit initialize configuration file
type OptInit struct {
	Init     bool   `docopt:"init"`
	Repo     bool   `docopt:"repo"`
	Template bool   `docopt:"template"`
	Name     string `docopt:"<name>"`
	Path     string `docopt:"<path>"`
}

// TODO: Unit tests for Init

// Init install a template
func Init(args []string, inject *di.DI) int {
	opts := &OptInit{}
	options.Bind(UsageDoc, args, opts)
	inst := inject.MakeInit()
	switch {
	case opts.Repo:
		return inst.CreateRepo(opts.Name, opts.Path)
	case opts.Template:
		return inst.CreateTemplate(opts.Name, opts.Path)
	}
	errs.Panic(errors.New(`Unknown error creating repo`))
	return 255
}
