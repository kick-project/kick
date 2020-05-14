package queries

import (
	"fmt"
	"testing"

	"github.com/crosseyed/prjstart/internal/db/dbinit"
)

func TestMatchTemplate(t *testing.T) {
	p := dbinit.InitTestDB()
	q := New("sqlite3", p)

	d := q.DB
	sql := d.Open()
	f := func(query string, args ...interface{}) {
		_, err := sql.Exec(query, args...)
		if err != nil {
			t.Fatal(err)
		}
	}

	f("INSERT OR REPLACE INTO global (name, url, desc) VALUES (?, ?, ?)", "publicglobal", "https://github.com/crosseyed/publicglobal", "Global Public")
	f("INSERT OR REPLACE INTO master (globalid, name, url, desc) SELECT id, ?, ?, ? FROM global WHERE url=?", "publicmaster", "https://github.com/crosseyed/crosseyed/master1", "Master 1", "https://github.com/crosseyed/publicglobal")
	f("INSERT OR REPLACE INTO templates (masterid, name, url, desc) SELECT id, ?, ?, ? FROM master WHERE url=?", "template1", "https://github.com/crosseyed/crosseyed/template1.git", "Template 1", "https://github.com/crosseyed/crosseyed/master1")

	fnrow := func(name, url, desc string, master_name, master_url, master_desc, global_name, global_url, global_desc string) {
		fmt.Printf("----\nname=%s\nurl=%s\ndescription=%s\nmaster_name=%s\nglobal_name=%s\n\n", name, url, desc, master_name, global_name)
	}
	q.SearchTemplate("template1", fnrow)
}
