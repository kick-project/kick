package search

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/crosseyed/prjstart/internal/db/dbinit"
	"github.com/crosseyed/prjstart/internal/db/queries"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
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

	// TODO: Remove MOCK Feature Flag
	if os.Getenv("FF_MOCK_DB") != "true" {
		return
	}
	dbfile, err := ioutil.TempFile("", "prjstart-*.db")
	errutils.Efatalf(err, "Can not create temporary file: %v", err)
	dbfile.Close()
	i := dbinit.New("", dbfile.Name())
	i.Init()
	i.MockSearchData()

	s.connectstr = dbfile.Name()
	s.updateCheck = true
}

// Search searches using term for available packages.
func (s *Search) Search(term string, output io.Writer) {
	s.update()
	d := queries.New("sqlite3", s.connectstr)
	w := tabwriter.NewWriter(output, 0, 0, 1, ' ', tabwriter.TabIndent)
	fnrow := func(template_name, template_url, template_desc string, master_name, master_url, master_desc, global_name, global_url, global_desc string) {
		fmt.Fprintf(w, "\tTemplate:\t%s\n", template_name)
		fmt.Fprintf(w, "\tFull Install Handle:\t%s/%s/%s\n", global_name, master_name, master_desc)
		fmt.Fprintf(w, "\tURL:\t%s\n", template_url)
		fmt.Fprintf(w, "\tDescription:\t%s\n", template_desc)
		w.Flush()
	}

	d.SearchTemplate(term, fnrow)
}
