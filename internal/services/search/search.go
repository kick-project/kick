// Package search implements search functionality
package search

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/services/search/entry"
	"github.com/kick-project/kick/internal/services/search/formatter"
	"gorm.io/gorm"
)

var querySearch = `
SELECT templateName, templateURL, templateDesc, repoName, repoURL, repoDesc
FROM
(
	SELECT
		template.name AS templateName,
		template.url AS templateURL,
		template.desc AS templateDesc,
		repo.name AS repoName,
		repo.url AS repoURL,
		repo.desc AS repoDesc,
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
		CASE WHEN LOWER(repo.name) LIKE LOWER(?)
		THEN
			template.name
		ELSE
			NULL
		END AS match3
	FROM template LEFT JOIN repo_template ON (template.id = repo_template.template_id)
	LEFT JOIN repo ON (repo_template.repo_id = repo.id)
	WHERE match1 IS NOT NULL OR match2 IS NOT NULL OR match3 IS NOT NULL
	ORDER BY
		match1 ASC NULLS LAST,
		match2 ASC NULLS LAST,
		match3 ASC NULLS LAST
)
`

// Search search for templates
type Search struct {
	Format formatter.Format
	ORM    *gorm.DB  `copier:"must"`
	Writer io.Writer `copier:"must"`
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
		errs.PanicF("query error: %w", err)
		defer rows.Close()

		for rows.Next() {
			var (
				name     sql.NullString
				URL      sql.NullString
				desc     sql.NullString
				repoName sql.NullString
				repoURL  sql.NullString
				repoDesc sql.NullString
			)
			err := rows.Scan(
				&name, &URL, &desc,
				&repoName, &repoURL, &repoDesc,
			)
			errs.FatalF("%v", err)

			ch <- &entry.Entry{
				Name:     name.String,
				URL:      URL.String,
				Desc:     desc.String,
				RepoName: repoName.String,
				RepoURL:  repoURL.String,
				RepoDesc: repoDesc.String,
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
