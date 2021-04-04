package isearch

import (
	"database/sql"
	"io"
	"os"

	"github.com/kick-project/kick/internal/services/search/formatter"
	"github.com/kick-project/kick/internal/settings"
	"gorm.io/gorm"
)

// Inject creates settings for search.New
func Inject(s *settings.Settings) (opts struct {
	DB     *sql.DB
	ORM    *gorm.DB
	Format formatter.Format
	Writer io.Writer
}) {
	format := &formatter.Standard{
		NoANSICodes: s.NoColour,
	}
	opts.DB = s.GetDB()
	opts.ORM = s.GetORM()
	opts.Format = format
	opts.Writer = os.Stdout
	return opts
}
