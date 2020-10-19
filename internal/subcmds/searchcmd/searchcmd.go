package searchcmd

import (
	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/services/search"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/isearch"
	"github.com/kick-project/kick/internal/settings/isync"
	"github.com/kick-project/kick/internal/utils/options"
	"github.com/jinzhu/copier"
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
	copier.Copy(synchro, isync.Inject(s))
	synchro.Templates()
	srch := &search.Search{}
	copier.Copy(srch, isearch.Inject(s))
	return srch.Search2Output(opts.Term)
}
