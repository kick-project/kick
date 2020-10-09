package internal

import (
	"github.com/crosseyed/prjstart/internal/fflags"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/docopt/docopt-go"
)

//
// Document strings
//

var UsageMain = `Generate project scaffolding from a predefined set of templates

Usage:
    prjstart start
    prjstart init
    prjstart list

Options:
    -h --help     Print help.
    -v --version  Print version.
    start         Start a project.
    init          Initialize configuration.
    list          List projects/variables.
`
var UsageMainFFREMOTE = `Generate project scaffolding from a predefined set of templates

Usage:
    prjstart start
    prjstart init
    prjstart install
    prjstart list
    prjstart search
    prjstart update

Options:
    -h --help     Print help.
    -v --version  Print version.
    start         Start a project.
    install       Install a template.
    update        Update repository data.
    init          Initialize configuration.
    list          List available project options.
    search        Search for available templates.
`

//
// Options
//

// OptMain holds all parsed options from GetOptMain.
type OptMain struct {
	Start   bool `docopt:"start"`
	Init    bool `docopt:"init"`
	Install bool `docopt:"install"`
	List    bool `docopt:"list"`
	Search  bool `docopt:"search"`
	Update  bool `docopt:"update"`
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
	if fflags.Remote() {
		opts, err = docopt.ParseArgs(UsageMainFFREMOTE, filterArgs, Version)
	} else {
		opts, err = docopt.ParseArgs(UsageMain, filterArgs, Version)
	}
	errutils.Epanicf("Can not parse usage doc: %s", err) // nolint
	o := new(OptMain)
	err = opts.Bind(o)
	errutils.Epanicf("Can not bind to structure: %s", err) // nolint
	return o
}
