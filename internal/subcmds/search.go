package subcmds

import (
	"github.com/crosseyed/prjstart/internal"
)

// Search searches for available templates
func Search(args []string) int {
	opts := internal.GetOptSearch(args)
	return SearchShort(opts)
}

// SearchShort prints the short output
func SearchShort(opts *internal.OptSearch) int {
	return 0
}
