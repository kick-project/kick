package updatecmd

import (
	"fmt"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `update repository data

Usage:
    kick update

Options:
    -h --help     print help
`

// OptUpdate bindings for docopts
type OptUpdate struct {
	Search bool `docopt:"update"`
}

// Update for templates
func Update(args []string, inject *di.DI) int {
	opts := &OptUpdate{}
	options.Bind(UsageDoc, args, opts)

	chk := inject.MakeCheck()

	if err := chk.Init(); err != nil {
		fmt.Fprintf(inject.Stderr, "%s\n", err.Error())
		exit.Exit(255)
	}

	u := inject.MakeUpdate()
	err := u.Build()
	errutils.Epanic(err)

	return 0
}
