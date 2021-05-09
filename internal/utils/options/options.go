package options

import (
	"github.com/docopt/docopt-go"
	"github.com/kick-project/kick/internal/resources/errs"
)

// Bind parse options from document text and populate a struct "opts".
// See https://pkg.go.dev/github.com/docopt/docopt-go#Opts.Bind for more
// information.
func Bind(usage string, args []string, opts interface{}) {
	parser, err := docopt.ParseArgs(usage, args, "")
	errs.PanicF("Can not parse usage doc: %s", err) // nolint
	err = parser.Bind(opts)
	errs.PanicF("Can not bind to structure: %s", err) // nolint
}
