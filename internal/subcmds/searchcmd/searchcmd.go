package searchcmd

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/services/search"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/icheck"
	"github.com/kick-project/kick/internal/settings/isearch"
	"github.com/kick-project/kick/internal/settings/isync"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `Search for templates using a keyword

Usage:
    kick search [-l] <term>

Options:
    -h --help  Print help.
    -l         Long output.
`

// OptSearch bindings for docopts
type OptSearch struct {
	Search bool   `docopt:"search"`
	Term   string `docopt:"<term>"`
	Long   bool   `docopt:"-l"`
}

// Search for templates
func Search(args []string, s *settings.Settings) int {
	opts := &OptSearch{}
	options.Bind(usageDoc, args, opts)

	chk := &check.Check{}
	err := copier.Copy(chk, icheck.Inject(s))
	errutils.Epanic(err)

	if err = chk.Init(); err != nil {
		fmt.Fprintf(s.Stderr, "%s\n", err.Error())
		utils.Exit(255)
	}

	// Sync DB table "installed" with configuration file
	synchro := &sync.Sync{}
	err = copier.Copy(synchro, isync.Inject(s))
	errutils.Epanic(err)
	synchro.Templates()
	srch := &search.Search{}
	err = copier.Copy(srch, isearch.Inject(s))
	errutils.Epanic(err)
	return srch.Search2Output(opts.Term)
}
