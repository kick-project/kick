package search

import (
	"github.com/crosseyed/prjstart/internal/db/queries"
)

// Search searches for available matching templates
type Search struct {
	srcurl      string
	connectstr  string
	queries     *queries.SQL
	updateCheck bool
}

func (s *Search) update() {
	if s.updateCheck {
		return
	}
}

// Search searches using term for available packages.
func (s *Search) Search(term string, fnrow func(template_name, template_url, template_desc, master_name, master_url, master_desc, global_name, global_url, globalc_desc string)) {
}
