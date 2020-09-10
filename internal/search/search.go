package search

import (
	"io/ioutil"
	"os"

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
	errutils.Efatalf("Can not create temporary file: %v", err)
	dbfile.Close()
	i := dbinit.New("", dbfile.Name())
	i.Init()
	i.MockSearchData()

	s.connectstr = dbfile.Name()
	s.updateCheck = true
}

// Search searches using term for available packages.
func (s *Search) Search(term string, fnrow func(template_name, template_url, template_desc, master_name, master_url, master_desc, global_name, global_url, globalc_desc string)) {
	s.update()
	d := queries.New("sqlite3", s.connectstr)
	d.SearchTemplate(term, fnrow)
}
