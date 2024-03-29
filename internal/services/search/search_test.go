package search_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestSearch(t *testing.T) {
	// Initialize database
	home := filepath.Join(testtools.TempDir(), "TestSearch")
	inject := di.New(&di.Options{Home: home})
	i := inject.MakeSetup()
	if _, err := os.Stat(inject.SqliteDB); err == nil {
		err = os.Remove(inject.SqliteDB)
		if err != nil {
			t.Error(err)
		}
	}
	i.Init()
	db := inject.MakeORM()
	buildSearchDataORM(t, db)

	// Target search term
	searchTerm := "template"

	srch := inject.MakeSearch()

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
		t.Logf("%v", entry)
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

	totalTestRepo2Rows := 0
	for range srch.Search("testrepo2") {
		totalTestRepo2Rows++
	}

	totalBoilerplate2Rows := 0
	for range srch.Search("boilerplate") {
		totalBoilerplate2Rows++
	}

	// Match expected counts
	assert.Equal(t, 2, counts["namePrefix"])
	assert.Equal(t, 4, counts["nameContains"])
	assert.Equal(t, 5, counts["urlContains"])
	assert.Equal(t, 7, totalTemplateRows)
	assert.Equal(t, 5, totalTestRepo2Rows)
	assert.Equal(t, 4, totalBoilerplate2Rows)
}

//
// Utils
//
func buildSearchDataORM(t *testing.T, db *gorm.DB) {
	m1 := model.Repo{
		Name: "testrepo",
		URL:  "http://127.0.0.1:8080/repo1.git",
		Desc: "Repo 1",
	}

	insertClause := clause.Insert{Modifier: "OR IGNORE"}

	result := db.Debug().Clauses(insertClause).Create(&m1)
	if result.Error != nil {
		t.Error(result.Error)
		return
	}

	t1 := []model.Template{
		{
			Name: "firsttemplate",
			URL:  "http://127.0.0.1:8080/tmpl.git",
			Desc: "My First Template",
		},
		{
			Name: "template1",
			URL:  "http://127.0.0.1:8080/tmpl1.git",
			Desc: "My Template Description",
		},
		{
			Name: "template2",
			URL:  "http://127.0.0.1:8080/tmpl2.git",
			Desc: "My Template Description",
		},
		{
			Name: "boilerplate3",
			URL:  "http://127.0.0.1:8080/template3.git",
			Desc: "My Template Description",
		},
	}

	for _, template := range t1 {
		template.Repo = []model.Repo{m1}
		db.Debug().Clauses(insertClause).Create(&template)
		if result.Error != nil {
			t.Error(result.Error)
			return
		}
	}

	m2 := model.Repo{
		Name: "testrepo2",
		URL:  "http://127.0.0.1:8080/repo2.git",
		Desc: "Repo 2",
	}

	result = db.Debug().Clauses(insertClause).Create(&m2)
	if result.Error != nil {
		t.Error(result)
		return
	}

	t2 := []model.Template{
		{
			Name: "mytemplate1",
			URL:  "http://127.0.0.1:8080/tmpl3.git",
			Desc: "My Boiler Plate Description",
		},
		{
			Name: "mytemplate1",
			URL:  "http://127.0.0.1:8080/tmpl4.git",
			Desc: "My Boiler Plate Description",
		},
		{
			Name: "mytemplate1",
			URL:  "http://127.0.0.1:8080/tmpl4.git",
			Desc: "My Boiler Plate Description",
		},
		{
			Name: "boilerplate1",
			URL:  "http://127.0.0.1:8080/template1.git",
			Desc: "My Boiler Plate Description",
		},
		{
			Name: "boilerplate2",
			URL:  "http://127.0.0.1:8080/template2.git",
			Desc: "My Template Description",
		},
		{
			Name: "boilerplate3",
			URL:  "http://127.0.0.1:8080/boilerplate4.git",
			Desc: "My Template Description",
		},
		{
			Name: "boilerplate4",
			URL:  "http://127.0.0.1:8080/boilerplate4.git",
			Desc: "My Template Description",
		},
	}

	for _, template := range t2 {
		template.Repo = []model.Repo{m2}
		db.Debug().Clauses(insertClause).Create(&template)
		if result.Error != nil {
			t.Error(result.Error)
			return
		}
	}

	t3 := []model.Template{
		{
			Name: "mytemplate4",
			URL:  "http://127.0.0.1:8080/template4.git",
			Desc: "My Boiler Plate Description",
		},
		{
			Name: "mytemplate5",
			URL:  "http://127.0.0.1:8080/template5.git",
			Desc: "My Boiler Plate Description",
		},
	}

	for _, template := range t3 {
		db.Debug().Clauses(insertClause).Create(&template)
		if result.Error != nil {
			t.Error(result.Error)
			return
		}
	}
}
