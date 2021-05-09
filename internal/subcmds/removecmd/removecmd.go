package removecmd

import (
	"errors"
	"fmt"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `Remove an installed template

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

// Remove remove a template
func Remove(args []string, inject *di.DI) int {
	opts := &OptRemove{}
	options.Bind(usageDoc, args, opts)
	if !opts.Remove {
		errutils.Epanic(errors.New("Remove set to false"))
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
