package searchcmd

import (
	"fmt"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/options"
)

// UsageDoc help document passed to docopts
var UsageDoc = `search for templates using a keyword

Usage:
    kick search [-l] <term>

Options:
    -h --help  print help
    -l         long output
    <term>     search term
`

// OptSearch bindings for docopts
type OptSearch struct {
	Search bool   `docopt:"search"`
	Term   string `docopt:"<term>"`
	Long   bool   `docopt:"-l"`
}

// Search for templates
func Search(args []string, inject *di.DI) int {
	opts := &OptSearch{}
	options.Bind(UsageDoc, args, opts)

	chk := inject.MakeCheck()
	if err := chk.Init(); err != nil {
		fmt.Fprintf(inject.Stderr, "%s\n", err.Error())
		exit.Exit(255)
	}

	// Sync DB table "installed" with configuration file
	synchro := inject.MakeSync()
	synchro.Files()
	srch := inject.MakeSearch()
	return srch.Search2Output(opts.Term)
}
