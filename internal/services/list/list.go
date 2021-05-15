package list

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/kick-project/kick/internal/resources/config"
	"github.com/olekukonko/tablewriter"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

// List manage listing of installed templates
type List struct {
	Stderr io.Writer    `validate:"required"`
	Stdout io.Writer    `validate:"required"`
	Conf   *config.File `validate:"required"`
}

// List lists the output
func (l *List) List(long bool) int {
	if long {
		l.longFmt()
	} else {
		l.shortFmt()
	}
	return 0
}

func (l *List) shortFmt() {
	templates := l.sort(l.Conf.Templates)
	sort.Sort(config.SortByName(templates))
	data := []string{}
	for _, t := range templates {
		name := t.Handle
		data = append(data, name)
	}
	termwidth, _ := terminal.Width()
	w := tabwriter.NewWriter(l.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	outstring := ""
	for _, name := range data {
		l := uint(0)
		if len(outstring) != 0 {
			l = uint(len(outstring))
			l += uint(len(name))
			l++
		}
		switch {
		case l == 0:
			outstring += name
		case l <= termwidth:
			outstring += "\t" + name
		default:
			fmt.Fprintln(w, outstring)
			outstring = name
		}
	}
	fmt.Fprintln(w, outstring)
	w.Flush()
}

func (l *List) longFmt() {
	var (
		header []string
		table  [][]string
	)
	header = []string{"Handle", "Description", "Template", "Location"}
	for _, row := range l.sort(l.Conf.Templates) {
		var (
			templateName string
			desc         string
		)
		if len(row.Template) > 0 {
			templateName = row.Template
			if len(row.Origin) > 0 {
				templateName = templateName + "/" + row.Origin
			}
		} else {
			templateName = "-"
		}
		desc = row.Desc
		if desc == "" {
			desc = "-"
		}
		table = append(table, []string{row.Handle, desc, templateName, row.URL})
	}
	writer := tablewriter.NewWriter(l.Stdout)
	writer.SetAlignment(tablewriter.ALIGN_LEFT)
	writer.SetHeader(header)
	for _, v := range table {
		writer.Append(v)
	}
	writer.Render()
}

func (l *List) sort(in []config.Template) (out []config.Template) {
	sort.Sort(config.SortByName(in))
	out = append(out, in...)
	return
}
