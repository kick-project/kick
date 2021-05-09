package listcmd

import (
	"fmt"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/utils/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `list handles/templates

Usage:
    kick list [-l]

Options:
    -h --help     print help
    -l            print long output
`

// OptList docopts options to list installed templates
type OptList struct {
	List bool `docopt:"list"`
	Long bool `docopt:"-l"`
}

// List starts the list sub command
func List(args []string, inject *di.DI) int {
	opts := &OptList{}
	options.Bind(UsageDoc, args, opts)

	chk := inject.MakeCheck()
	if err := chk.Init(); err != nil {
		fmt.Fprintf(inject.Stderr, "%s\n", err.Error())
		exit.Exit(255)
	}

	l := inject.MakeList()
	return l.List(opts.Long)
}
