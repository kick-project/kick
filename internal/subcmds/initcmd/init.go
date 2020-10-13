package initcmd

import (
	"log"

	"github.com/crosseyed/prjstart/internal/services/initialize"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/settings/iinitialize"
	"github.com/crosseyed/prjstart/internal/utils/options"
)

var usageDoc = `Initialize configuration

Usage:
    prjstart init

Options:
    -h --help     Print help.
`

// OptInit initialize configuration file
type OptInit struct {
	Init bool `docopt:"init"`
}

// InitCmd initialize configuration
func InitCmd(args []string, s *settings.Settings) int {
	opts := &OptInit{}
	options.Bind(usageDoc, args, opts)
	if opts.Init == false {
		log.Println("error can not initialize")
		return 256
	}

	i := initialize.New(iinitialize.Inject(s))
	i.Init()

	return 0
}
