package initcmd

import (
	"log"

	"github.com/crosseyed/prjstart/internal/services/initialize"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/docopt/docopt-go"
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

// GetOptInit parse document for init options
func GetOptInit(args []string) *OptInit {
	opts, err := docopt.ParseArgs(usageDoc, args, "")
	errutils.Epanicf("Can not parse usage doc: %s", err) // nolint
	o := new(OptInit)
	err = opts.Bind(o)
	errutils.Epanicf("Can not bind to structure: %s", err) // nolint
	return o
}

// InitCmd initialize configuration
func InitCmd(args []string, s *settings.Settings) int {
	opts := GetOptInit(args)
	if opts.Init == false {
		log.Println("error can not initialize")
		return 256
	}

	i := initialize.New(s.Initialize())
	i.Init()
	return 0
}
