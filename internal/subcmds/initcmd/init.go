package initcmd

import (
	"log"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/iinitialize"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
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
	if !opts.Init {
		log.Println("error can not initialize")
		return 256
	}

	i := &initialize.Initialize{}
	err := copier.Copy(i, iinitialize.Inject(s))
	errutils.Epanic(err)
	i.Init()

	return 0
}
