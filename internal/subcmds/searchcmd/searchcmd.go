package searchcmd

import (
	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/services/search"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/isearch"
	"github.com/kick-project/kick/internal/settings/isync"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `Search for templates using a keyword

Usage:
    prjstart search [--long] <term>

Options:
    -h --help     Print help.
    -l --long     Long output.
`

// OptSearch bindings for docopts
type OptSearch struct {
	Search bool   `docopt:"search"`
	Term   string `docopt:"<term>"`
	Long   bool   `docopt:"--long"`
}

// Search for templates
func Search(args []string, s *settings.Settings) int {
	opts := &OptSearch{}
	options.Bind(usageDoc, args, opts)

	// Sync DB table "installed" with configuration file
	synchro := &sync.Sync{}
	err := copier.Copy(synchro, isync.Inject(s))
	errutils.Epanic(err)
	synchro.Templates()
	srch := &search.Search{}
	err = copier.Copy(srch, isearch.Inject(s))
	errutils.Epanic(err)
	return srch.Search2Output(opts.Term)
}
