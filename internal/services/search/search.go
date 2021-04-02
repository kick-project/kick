// Package search implements search functionality
package search

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/kick-project/kick/internal/services/search/entry"
	"github.com/kick-project/kick/internal/services/search/formatter"
	"github.com/kick-project/kick/internal/utils/errutils"
	"gorm.io/gorm"
)

var querySearch = `
SELECT templateName, templateURL, templateDesc, masterName, masterURL, masterDesc
FROM
(
	SELECT
		template.name AS templateName,
		template.url AS templateURL,
		template.desc AS templateDesc,
		master.name AS masterName,
		master.url AS masterURL,
		master.desc AS masterDesc,
		CASE WHEN LOWER(template.name) LIKE LOWER(?)
			THEN
				template.name
			ELSE
				NULL
		END AS match1,
		CASE WHEN LOWER(template.url) LIKE LOWER(?)
			THEN
				template.name
			ELSE
				NULL
		END AS match2,
		CASE WHEN LOWER(master.name) LIKE LOWER(?)
		THEN
			template.name
		ELSE
			NULL
		END AS match3
	FROM template LEFT JOIN master_template ON (template.id = master_template.template_id)
	LEFT JOIN master ON (master_template.master_id = master.id)
	WHERE match1 IS NOT NULL OR match2 IS NOT NULL OR match3 IS NOT NULL
	ORDER BY
		match1 ASC NULLS LAST,
		match2 ASC NULLS LAST,
		match3 ASC NULLS LAST
)
`

// Search search for templates
type Search struct {
	ORM    *gorm.DB
	Format formatter.Format
	Writer io.Writer
}

// Search searches database for term and returns the results through *Entry channel.
func (s *Search) Search(term string) <-chan *entry.Entry {
	ch := make(chan *entry.Entry, 24)
	go func() {
		var err error
		rows, err := s.ORM.Raw(
			querySearch,
			fmt.Sprintf("%s%%", term),
			fmt.Sprintf("%%%s%%", term),
			fmt.Sprintf("%%%s%%", term),
		).Rows()
		errutils.Epanicf("query error: %w", err)
		defer rows.Close()

		for rows.Next() {
			var (
				name       sql.NullString
				URL        sql.NullString
				desc       sql.NullString
				masterName sql.NullString
				masterURL  sql.NullString
				masterDesc sql.NullString
			)
			err := rows.Scan(
				&name, &URL, &desc,
				&masterName, &masterURL, &masterDesc,
			)
			errutils.Efatalf("%v", err)

			ch <- &entry.Entry{
				Name:       name.String,
				URL:        URL.String,
				Desc:       desc.String,
				MasterName: masterName.String,
				MasterURL:  masterURL.String,
				MasterDesc: masterDesc.String,
			}
		}
		close(ch)
	}()
	return ch
}

// Search2Output searches database for term and sends the results to the formatter.Format function supplied in New.
// Blocks until all entries are processed.
func (s *Search) Search2Output(term string) int {
	ch := s.Search(term)
	if s.Format != nil {
		s.Format.Writer(s.Writer, ch)
	} else {
		fmtter := formatter.Standard{}
		fmtter.Writer(s.Writer, ch)
	}
	return 0
}
