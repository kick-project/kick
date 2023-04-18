package templatescan

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/kick-project/kick/internal/resources/model"
	"github.com/kick-project/kick/internal/resources/testtools"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestScan_Run(t *testing.T) {
	root := filepath.Join(testtools.FixtureDir(), "gotemplate")
	dbfile := filepath.Join(testtools.TempDir(), "TestScan_Run.db")
	db := model.CreateModelTemporary(model.Options{File: dbfile})
	s := Scan{
		DB: db,
	}
	err := s.Run(root, 5)
	assert.NoError(t, err)
	testOptions(t, db)
	testLabels(t, db)
}

func testOptions(t *testing.T, db *gorm.DB) {
	q := QueryScanOption + " WHERE file = ?"
	type Result struct {
		Dir    string
		Path   string
		Option string
	}
	results := []Result{}

	db.Raw(q, "go.mod").Scan(&results)
	assert.Len(t, results, 1)
	assert.Equal(t, "render", results[0].Option)
}

func testLabels(t *testing.T, db *gorm.DB) {
	q := QueryScanLabel + " WHERE label = ?"
	type Result struct {
		Dir   string
		Path  string
		Label string
	}
	results := []Result{}

	db.Raw(q, "editor").Scan(&results)
	assert.Len(t, results, 1)
	assert.Equal(t, ".editorconfig", results[0].Path)

	db.Raw(q, "go").Scan(&results)
	assert.Len(t, results, 1)
	assert.Equal(t, "go.mod", results[0].Path)

	db.Raw(q, "github").Scan(&results)
	assert.Len(t, results, 3)
	for _, r := range results {
		assert.True(t, strings.HasPrefix(r.Path, ".github"))
	}
}
