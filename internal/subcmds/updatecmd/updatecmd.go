package updatecmd

import (
	"github.com/kick-project/kick/internal/services/update"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/iupdate"
	"github.com/kick-project/kick/internal/utils/options"
	"github.com/jinzhu/copier"
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
	copier.Copy(u, iupdate.Inject(s))
	u.Build()

	return 0
}
