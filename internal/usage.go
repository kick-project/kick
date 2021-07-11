package internal

import (
	"github.com/docopt/docopt-go"
	"github.com/kick-project/kick/internal/resources/errs"
)

//
// Document strings
//

var usageDoc = `Generate project scaffolding from a predefined set of templates

Usage:
    kick start
    kick install
    kick remove
    kick search
    kick update
    kick setup
    kick init
    kick repo

Options:
    -h --help     print help
    -v --version  print version
    start         start a project
    install       install a template
    remove        remove an installed template
    search        search repositories for available templates
    update        update local repository information 
    setup         setup configuration
    init          initialize a template or repository
    repo          tool to build a repository
`

//
// Options
//

// OptMain holds all parsed options from GetOptMain.
type OptMain struct {
	Start   bool `docopt:"start"`
	Setup   bool `docopt:"setup"`
	Install bool `docopt:"install"`
	List    bool `docopt:"list"`
	Remove  bool `docopt:"remove"`
	Search  bool `docopt:"search"`
	Update  bool `docopt:"update"`
	Init    bool `docopt:"init"`
	Repo    bool `docopt:"repo"`
}

// GetOptMain is a command line option parser that uses docopts-go to parse a
// usage document string.
func GetOptMain(args []string) *OptMain {
	var (
		opts docopt.Opts
		err  error
	)
	filterArgs := []string{}
	i := 0
	for _, arg := range args {
		i++
		if i == 1 {
			continue
		}
		filterArgs = append(filterArgs, arg)
		break
	}
	opts, err = docopt.ParseArgs(usageDoc, filterArgs, Version)
	errs.PanicF("Can not parse usage doc: %s", err) // nolint
	o := new(OptMain)
	err = opts.Bind(o)
	errs.PanicF("Can not bind to structure: %s", err) // nolint
	return o
}
