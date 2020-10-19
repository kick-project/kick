package updatecmd

import (
	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/services/update"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/iupdate"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
)

var usageDoc = `Update repository data

Usage:
    prjstart update

Options:
    -h --help     Print help.
`

// OptUpdate bindings for docopts
type OptUpdate struct {
	Search bool `docopt:"update"`
}

// Update for templates
func Update(args []string, s *settings.Settings) int {
	opts := &OptUpdate{}
	options.Bind(usageDoc, args, opts)

	u := &update.Update{}
	err := copier.Copy(u, iupdate.Inject(s))
	errutils.Epanic(err)
	err = u.Build()
	errutils.Epanic(err)

	return 0
}
