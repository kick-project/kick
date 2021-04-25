package search

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/di/iinitialize"
	"github.com/kick-project/kick/internal/di/isearch"
	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/services/initialize"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestSearch(t *testing.T) {
	// Initialize database
	home := filepath.Join(utils.TempDir(), "TestSearch")
	inject := di.Setup(home)
	i := &initialize.Initialize{}
	err := copier.Copy(i, iinitialize.Inject(inject))
	errutils.Epanic(err)
	if _, err = os.Stat(inject.SqliteDB); err == nil {
		err = os.Remove(inject.SqliteDB)
		if err != nil {
			t.Error(err)
		}
	}
	i.Init()
	db := inject.GetORM()
	buildSearchDataORM(t, db)

	// Target search term
	searchTerm := "template"

	srch := &Search{}
	err = copier.Copy(srch, isearch.Inject(inject))
	errutils.Epanic(err)

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
	assert.Equal(t, 4, counts["nameContains"])
	assert.Equal(t, 5, counts["urlContains"])
	assert.Equal(t, 7, totalTemplateRows)
	assert.Equal(t, 5, totalTestMaster2Rows)
	assert.Equal(t, 4, totalBoilerplate2Rows)
}

//
// Utils
//
func buildSearchDataORM(t *testing.T, db *gorm.DB) {
	m1 := model.Master{
		Name: "testmaster",
		URL:  "http://127.0.0.1:5000/master1.git",
		Desc: "Master 1",
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
			URL:  "http://127.0.0.1:5000/tmpl.git",
			Desc: "My First Template",
		},
		{
			Name: "template1",
			URL:  "http://127.0.0.1:5000/tmpl1.git",
			Desc: "My Template Description",
		},
		{
			Name: "template2",
			URL:  "http://127.0.0.1:5000/tmpl2.git",
			Desc: "My Template Description",
		},
		{
			Name: "boilerplate3",
			URL:  "http://127.0.0.1:5000/template3.git",
			Desc: "My Template Description",
		},
	}

	for _, template := range t1 {
		template.Master = []model.Master{m1}
		db.Debug().Clauses(insertClause).Create(&template)
		if result.Error != nil {
			t.Error(result.Error)
			return
		}
	}

	m2 := model.Master{
		Name: "testmaster2",
		URL:  "http://127.0.0.1:5000/master2.git",
		Desc: "Master 2",
	}

	result = db.Debug().Clauses(insertClause).Create(&m2)
	if result.Error != nil {
		t.Error(result)
		return
	}

	t2 := []model.Template{
		{
			Name: "mytemplate1",
			URL:  "http://127.0.0.1:5000/tmpl3.git",
			Desc: "My Boiler Plate Description",
		},
		{
			Name: "mytemplate1",
			URL:  "http://127.0.0.1:5000/tmpl4.git",
			Desc: "My Boiler Plate Description",
		},
		{
			Name: "mytemplate1",
			URL:  "http://127.0.0.1:5000/tmpl4.git",
			Desc: "My Boiler Plate Description",
		},
		{
			Name: "boilerplate1",
			URL:  "http://127.0.0.1:5000/template1.git",
			Desc: "My Boiler Plate Description",
		},
		{
			Name: "boilerplate2",
			URL:  "http://127.0.0.1:5000/template2.git",
			Desc: "My Template Description",
		},
		{
			Name: "boilerplate3",
			URL:  "http://127.0.0.1:5000/boilerplate4.git",
			Desc: "My Template Description",
		},
		{
			Name: "boilerplate4",
			URL:  "http://127.0.0.1:5000/boilerplate4.git",
			Desc: "My Template Description",
		},
	}

	for _, template := range t2 {
		template.Master = []model.Master{m2}
		db.Debug().Clauses(insertClause).Create(&template)
		if result.Error != nil {
			t.Error(result.Error)
			return
		}
	}

	t3 := []model.Template{
		{
			Name: "mytemplate4",
			URL:  "http://127.0.0.1:5000/template4.git",
			Desc: "My Boiler Plate Description",
		},
		{
			Name: "mytemplate5",
			URL:  "http://127.0.0.1:5000/template5.git",
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
