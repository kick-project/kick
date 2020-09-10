package search

import (
	"bytes"
	"fmt"
	"testing"
	"text/tabwriter"
)

func TestSearch(t *testing.T) {
	buf := bytes.Buffer{}
	s := Search{}
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.TabIndent)
	fnrow := func(template_name, template_url, template_desc, master_name, master_url, master_desc, global_name, global_url, global_desc string) {
		fmt.Fprintf(w, template_name)
		w.Flush()
	}
	s.Search("template1", fnrow)
	if len(buf.String()) == 0 {
		t.Fail()
	}
}
