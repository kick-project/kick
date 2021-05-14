// DO NOT EDIT: Generated using "make interfaces"

package search

import (
	"github.com/kick-project/kick/internal/services/search/entry"
)

// SearchIface ...
type SearchIface interface {
	// Search searches database for term and returns the results through *Entry channel.
	Search(term string) <-chan *entry.Entry
	// Search2Output searches database for term and sends the results to the formatter.Format function supplied in New.
	// Blocks until all entries are processed.
	Search2Output(term string) int
}
