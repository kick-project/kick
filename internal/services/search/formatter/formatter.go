package formatter

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"

	"github.com/kick-project/kick/internal/services/search/entry"
)

// Format Format is a interface to format search entries.
type Format interface {
	// Writer Writer is abstract type that takes a search.Entry then writes to io.Writer
	Writer(io.Writer, <-chan *entry.Entry)
}

// Tables format as a table
type Tables struct {
	long bool
}

// New return Tables pointer
func New(long bool) *Tables {
	return &Tables{
		long: long,
	}
}

// Writer write entries
func (t *Tables) Writer(w io.Writer, ch <-chan *entry.Entry) {
	var (
		header []string
		table  [][]string
	)
	if t.long {
		header, table = t.longFmt(ch)
	} else {
		header, table = t.shortFmt(ch)
	}
	writer := tablewriter.NewWriter(w)
	writer.SetAlignment(tablewriter.ALIGN_LEFT)
	writer.SetHeader(header)
	for _, v := range table {
		writer.Append(v)
	}
	writer.Render()
}

func (t *Tables) longFmt(ch <-chan *entry.Entry) (header []string, table [][]string) {
	header = []string{"Template", "Repository", "Template description", "Template Location"}
	for e := range ch {
		row := []string{fmt.Sprintf("%s/%s", e.Name, e.RepoName), e.RepoName, e.Desc, e.URL}
		table = append(table, row)
	}
	return
}

func (t *Tables) shortFmt(ch <-chan *entry.Entry) (header []string, table [][]string) {
	header = []string{"Template", "Location"}
	for e := range ch {
		row := []string{fmt.Sprintf("%s/%s", e.Name, e.RepoName), e.URL}
		table = append(table, row)
	}
	return
}
