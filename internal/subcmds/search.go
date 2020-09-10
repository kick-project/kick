package subcmds

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/crosseyed/prjstart/internal"
	"github.com/crosseyed/prjstart/internal/search"
)

// Search searches for available templates
func Search(args []string) int {
	opts := internal.GetOptSearch(args)
	if opts.Long {
		return SearchLong(opts)
	}
	return SearchShort(opts)
}

// SearchShort prints the short output of the search results
func SearchShort(opts *internal.OptSearch) int {
	s := search.Search{}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fnrow := func(template_name, template_url, template_desc, master_name, master_url, master_desc, global_name, global_url, global_desc string) {
		fmt.Fprintf(w, "\tTemplate:\t%s\n", template_name)
		fmt.Fprintf(w, "\tFull Install Handle:\t%s/%s/%s\n\n", global_name, master_name, template_name)
		w.Flush()
	}
	s.Search(opts.Template, fnrow)
	return 0
}

// SearchLong prints the long output of the search results
func SearchLong(opts *internal.OptSearch) int {
	s := search.Search{}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	fnrow := func(template_name, template_url, template_desc, master_name, master_url, master_desc, global_name, global_url, global_desc string) {
		fmt.Fprintf(w, "\tTemplate:\t%s\n", template_name)
		fmt.Fprintf(w, "\tFull Install Handle:\t%s/%s/%s\n", global_name, master_name, template_desc)
		fmt.Fprintf(w, "\tURL:\t%s\n", template_url)
		fmt.Fprintf(w, "\tDescription:\t%s\n\n", template_desc)
		w.Flush()
	}
	s.Search(opts.Template, fnrow)
	return 0
}
