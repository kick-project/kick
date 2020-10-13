package search

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/crosseyed/prjstart/internal/resources/db"
	"github.com/crosseyed/prjstart/internal/services/initialize"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/settings/iinitialize"
	"github.com/crosseyed/prjstart/internal/settings/isearch"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

var insertMaster string = `---
INSERT OR IGNORE INTO master (name, url, desc) VALUES (?, ?, ?)
`

var insertTemplate string = `---
INSERT OR IGNORE INTO templates (masterid, name, url, desc)
	SELECT master.id, ?, ?, ? FROM master WHERE master.url = ?
`

func TestSearch(t *testing.T) {
	// Initialize database
	home := filepath.Join(utils.TempDir(), "TestSearch")
	s := settings.GetSettings(home)
	i := &initialize.Initialize{}
	copier.Copy(i, iinitialize.Inject(s))
	i.Init()
	dbconn := s.GetDB()
	buildSearchData(t, dbconn)

	// Target search term
	searchTerm := "template"

	srch := &Search{}
	copier.Copy(srch, isearch.Inject(s))

	// Parallel arrays to count regexp
	matchmaker := map[string]*regexp.Regexp{}
	counts := map[string]int{}

	rePrefix := regexp.MustCompile(fmt.Sprintf("(?i)^%s", searchTerm))
	reContains := regexp.MustCompile(fmt.Sprintf("(?i)%s", searchTerm))
	matchmaker["namePrefix"] = rePrefix
	matchmaker["nameContains"] = reContains
	matchmaker["urlContains"] = reContains
	for key := range matchmaker {
		counts[key] = 0
	}

	totalTemplateRows := 0
	for entry := range srch.Search("template") {
		totalTemplateRows++
		if matchmaker["namePrefix"].MatchString(entry.Name) {
			counts["namePrefix"]++
		}
		if matchmaker["nameContains"].MatchString(entry.Name) {
			counts["nameContains"]++
		}
		if matchmaker["urlContains"].MatchString(entry.URL) {
			counts["urlContains"]++
		}
	}

	totalTestMaster2Rows := 0
	for range srch.Search("testmaster2") {
		totalTestMaster2Rows++
	}

	totalBoilerplate2Rows := 0
	for range srch.Search("boilerplate") {
		totalBoilerplate2Rows++
	}

	// Match expected counts
	assert.Equal(t, 2, counts["namePrefix"])
	assert.Equal(t, 2, counts["nameContains"])
	assert.Equal(t, 3, counts["urlContains"])
	assert.Equal(t, 5, totalTemplateRows)
	assert.Equal(t, 5, totalTestMaster2Rows)
	assert.Equal(t, 4, totalBoilerplate2Rows)
}

//
// Utils
//

func buildSearchData(t *testing.T, dbconn *sql.DB) {
	masterURL1 := "http://127.0.0.1:5000/master1.git"
	masterURL2 := "http://127.0.0.1:5000/master2.git"
	masters := []map[string]string{
		{
			"Name": "testmaster",
			"URL":  masterURL1,
			"Desc": "My Master Description",
		},
		{
			"Name": "testmaster2",
			"URL":  masterURL2,
			"Desc": "My Master Description",
		},
	}
	templates1 := []map[string]string{
		{
			"Name": "firsttemplate",
			"URL":  "http://127.0.0.1:5000/tmpl.git",
			"Desc": "My First Template",
		},
		{
			"Name": "template1",
			"URL":  "http://127.0.0.1:5000/tmpl1.git",
			"Desc": "My Template Description",
		},
		{
			"Name": "template2",
			"URL":  "http://127.0.0.1:5000/tmpl2.git",
			"Desc": "My Template Description",
		},
		{
			"Name": "boilerplate3",
			"URL":  "http://127.0.0.1:5000/template3.git",
			"Desc": "My Template Description",
		},
	}
	templates2 := []map[string]string{
		{
			"Name": "mytemplate1",
			"URL":  "http://127.0.0.1:5000/tmpl3.git",
			"Desc": "My Boiler Plate Description",
		},
		{
			"Name": "mytemplate1",
			"URL":  "http://127.0.0.1:5000/tmpl4.git",
			"Desc": "My Boiler Plate Description",
		},
		{
			"Name": "mytemplate1",
			"URL":  "http://127.0.0.1:5000/tmpl4.git",
			"Desc": "My Boiler Plate Description",
		},
		{
			"Name": "boilerplate1",
			"URL":  "http://127.0.0.1:5000/template1.git",
			"Desc": "My Boiler Plate Description",
		},
		{
			"Name": "boilerplate2",
			"URL":  "http://127.0.0.1:5000/template2.git",
			"Desc": "My Template Description",
		},
		{
			"Name": "boilerplate3",
			"URL":  "http://127.0.0.1:5000/boilerplate4.git",
			"Desc": "My Template Description",
		},
		{
			"Name": "boilerplate4",
			"URL":  "http://127.0.0.1:5000/boilerplate4.git",
			"Desc": "My Template Description",
		},
	}

	db.Lock()
	defer db.Unlock()

	for _, m := range masters {
		_, err := dbconn.Exec(insertMaster, m["Name"], m["URL"], m["Desc"])
		if err != nil {
			t.Errorf("%w", err)
		}
	}

	for _, tpl := range templates1 {
		_, err := dbconn.Exec(insertTemplate, tpl["Name"], tpl["URL"], tpl["Desc"], masterURL1)
		if err != nil {
			t.Errorf("%w", err)
		}
	}

	for _, tpl := range templates2 {
		_, err := dbconn.Exec(insertTemplate, tpl["Name"], tpl["URL"], tpl["Desc"], masterURL2)
		if err != nil {
			t.Errorf("%w", err)
		}
	}
}
