// Package search implements search functionality
package search

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/crosseyed/prjstart/internal/services/search/entry"
	"github.com/crosseyed/prjstart/internal/services/search/formatter"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

var querySearch = `
SELECT templateName, templateURL, templateDesc, masterName, masterURL, masterDesc
FROM
(
	SELECT
		templates.name AS templateName,
		templates.url AS templateURL,
		templates.desc AS templateDesc,
		master.name AS masterName,
		master.url AS masterURL,
		master.desc AS masterDesc,
		CASE WHEN LOWER(templates.name) LIKE LOWER(?)
			THEN
				templates.name
			ELSE
				NULL
		END AS match1,
		CASE WHEN LOWER(templates.url) LIKE LOWER(?)
			THEN
				templates.name
			ELSE
				NULL
		END AS match2,
		CASE WHEN LOWER(master.name) LIKE LOWER(?)
		THEN
			templates.name
		ELSE
			NULL
		END AS match3
	FROM templates LEFT JOIN master ON (templates.masterid = master.id)
	WHERE match1 IS NOT NULL OR match2 IS NOT NULL OR match3 IS NOT NULL
	ORDER BY
		match1 ASC NULLS LAST,
		match2 ASC NULLS LAST,
		match3 ASC NULLS LAST
)
`

// Search search for templates
type Search struct {
	db     *sql.DB
	format formatter.Format
	writer io.Writer
}

// Options options for New
type Options struct {
	DB     *sql.DB
	Format formatter.Format
	Writer io.Writer
}

// New creates a Search
func New(opts Options) *Search {
	if opts.DB == nil {
		panic("opts.DB is nil")
	}
	if opts.Writer == nil {
		panic("opts.Writer is nil")
	}
	s := &Search{
		db:     opts.DB,
		format: opts.Format,
		writer: opts.Writer,
	}
	return s
}

// Search searches database for term and returns the results through *Entry channel.
func (s *Search) Search(term string) <-chan *entry.Entry {
	ch := make(chan *entry.Entry, 24)
	go func() {
		var err error
		rows, err := s.db.Query(
			querySearch,
			fmt.Sprintf("%s%%", term),
			fmt.Sprintf("%%%s%%", term),
			fmt.Sprintf("%%%s%%", term),
		)
		errutils.Epanicf("query error: %w", err)
		defer rows.Close()

		for rows.Next() {
			curEntry := &entry.Entry{}
			err := rows.Scan(
				&curEntry.Name, &curEntry.URL, &curEntry.Desc,
				&curEntry.MasterName, &curEntry.MasterURL, &curEntry.MasterDesc,
			)
			errutils.Efatalf("%v", err)
			ch <- curEntry
		}
		close(ch)
	}()
	return ch
}

// Search2Output searches database for term and sends the results to the formatter.Format function supplied in New.
// Blocks until all entries are processed.
func (s *Search) Search2Output(term string) int {
	ch := s.Search(term)
	if s.format != nil {
		s.format.Writer(s.writer, ch)
	} else {
		fmtter := formatter.Standard{}
		fmtter.Writer(s.writer, ch)
	}
	return 0
}
