package internal

import (
	"github.com/docopt/docopt-go"
	"github.com/kick-project/kick/internal/utils/errutils"
)

//
// Document strings
//

var usageDoc = `Generate project scaffolding from a predefined set of templates

Usage:
    kick start
    kick init
    kick install
    kick list
    kick remove
    kick search
    kick update

Options:
    -h --help     print help
    -v --version  print version
    start         start a project
    init          initialize configuration
    install       install a template
    list          list available project options
    remove        remove an installed template
    search        search for available templates
    update        update repository data
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
	Remove  bool `docopt:"remove"`
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
	opts, err = docopt.ParseArgs(usageDoc, filterArgs, Version)
	errutils.Epanicf("Can not parse usage doc: %s", err) // nolint
	o := new(OptMain)
	err = opts.Bind(o)
	errutils.Epanicf("Can not bind to structure: %s", err) // nolint
	return o
}
