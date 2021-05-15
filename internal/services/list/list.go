package list

import (
	"fmt"
	"io"

	"github.com/kick-project/kick/internal/resources/config"
	"github.com/olekukonko/tablewriter"
)

// List manage listing of installed templates
type List struct {
	Stderr io.Writer    `validate:"required"`
	Stdout io.Writer    `validate:"required"`
	Conf   *config.File `validate:"required"`
}

// List lists the output
func (l *List) List(long bool) int {
	var (
		header []string
		table  [][]string
	)
	if long {
		header, table = l.longFmt()
	} else {
		header, table = l.shortFmt()
	}
	writer := tablewriter.NewWriter(l.Stdout)
	writer.SetAlignment(tablewriter.ALIGN_LEFT)
	writer.SetHeader(header)
	for _, v := range table {
		writer.Append(v)
	}
	writer.Render()
	return 0
}

func (l *List) shortFmt() (header []string, table [][]string) {
	header = []string{"Handle", "Description"}
	for _, row := range l.Conf.Templates {
		table = append(table, []string{row.Handle, row.Desc})
	}
	return
}

func (l *List) longFmt() (header []string, table [][]string) {
	header = []string{"Handle", "Description", "Template", "Location"}
	for _, row := range l.Conf.Templates {
		table = append(table, []string{row.Handle, row.Desc, fmt.Sprintf("%s/%s", row.Template, row.Origin), row.URL})
	}
	return
}
