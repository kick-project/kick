package isearch

import (
	"io"
	"os"

	"github.com/kick-project/kick/internal/di"
	"github.com/kick-project/kick/internal/services/search/formatter"
	"gorm.io/gorm"
)

// Inject creates di for search.New
func Inject(s *di.DI) (opts struct {
	ORM    *gorm.DB
	Format formatter.Format
	Writer io.Writer
}) {
	format := &formatter.Standard{
		NoANSICodes: s.NoColour,
	}
	opts.ORM = s.MakeORM()
	opts.Format = format
	opts.Writer = os.Stdout
	return opts
}
