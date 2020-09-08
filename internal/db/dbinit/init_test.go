package dbinit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/crosseyed/prjstart/internal/db"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

func TestInit(t *testing.T) {
	p, err := filepath.Abs(filepath.Join("..", "..", "..", "tmp", "TestInit.db"))
	if err != nil {
		t.Fatalf("filepath.Abs %s error: %v", p, err)
		t.Fail()
		return
	}
	os.Remove(p)
	i := New(p)
	i.Init()

	db_ := db.New("sqlite3", p)
	db_.Lock()
	d := db_.Open()
	defer db_.Unlock()
	defer db_.Close()

	f := func(query string, args ...interface{}) {
		_, err = d.Exec(query, args...)
		if err != nil {
			t.Fatal(err)
		}
	}
	f("INSERT OR REPLACE INTO global (name, url, desc) VALUES (?, ?, ?)", "TESTGLOBAL", "URLGLOBAL", "Test global")
	f("INSERT OR REPLACE INTO master (globalid, name, url, desc) SELECT id, ?, ?, ? FROM global WHERE url=?", "TESTMASTER", "URLMASTER", "Test master", "URLGLOBAL")
	f("INSERT OR REPLACE INTO templates (masterid, name, url, desc) SELECT id, ?, ?, ? FROM master WHERE url=?", "TESTTEMPLATE", "URLTEMPLATE", "Test template", "URLMASTER")

	rows, err := d.Query(`SELECT count(*) as count, templates.url as template, master.url as master, global.url as global FROM templates JOIN master JOIN global WHERE template=? AND master=? AND global=?`, "URLTEMPLATE", "URLMASTER", "URLGLOBAL")
	if err != nil {
		t.Fatal(err)
	}

	var (
		template string
		master   string
		global   string
		count    int
	)
	for rows.Next() {
		err := rows.Scan(&count, &template, &master, &global)
		if errutils.Efatalf(err, "Can not scan row: %v", err) {
			break
		}
	}

	if count != 1 || template != "URLTEMPLATE" || master != "URLMASTER" || global != "URLGLOBAL" {
		t.Fail()
	}
}
