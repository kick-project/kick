package updatecmd

import (
	"github.com/crosseyed/prjstart/internal/services/metadata"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/settings/imetadata"
	"github.com/crosseyed/prjstart/internal/utils/options"
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

	m := metadata.New(imetadata.Inject(s))
	m.Build()

	// TODO: Create a real return code.
	return 0
}
