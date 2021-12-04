package removecmd

import (
	"errors"
	"fmt"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `Remove an installed template

Usage:
    kick remove <handle>

Options:
    -h --help        print help
    <handle>         handle to remove
`

// OptRemove remove an installed template
type OptRemove struct {
	Remove bool   `docopt:"remove"`
	Handle string `docopt:"<handle>"`
}

// TODO: Unit test for removecommand.go

// Remove remove a template
func Remove(args []string, inject *di.DI) int {
	opts := &OptRemove{}
	options.Bind(UsageDoc, args, opts)
	if !opts.Remove {
		errs.Panic(errors.New("Remove set to false"))
		return 256
	}

	chk := inject.MakeCheck()
	if err := chk.Init(); err != nil {
		fmt.Fprintf(inject.Stderr, "%s\n", err.Error())
		exit.Exit(255)
	}

	rm := inject.MakeRemove()
	return rm.Remove(opts.Handle)
}
