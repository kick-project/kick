package list

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/kick-project/kick/internal/resources/ansicodes"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/mattn/go-isatty"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

// List manage listing of installed templates
type List struct {
	Stderr io.Writer    `copier:"must"`
	Stdout io.Writer    `copier:"must"`
	Conf   *config.File `copier:"must"`
}

// List lists the output
func (l *List) List(long bool) int {
	if long {
		return l.long()
	}
	return l.short()
}

// short provides a short list of templates.
func (l *List) short() int {
	templates := l.Conf.Templates
	sort.Sort(config.SortByName(templates))
	data := []string{}
	for _, t := range templates {
		name := t.Handle
		data = append(data, name)
	}
	l.shortOutput(data)
	return 0
}

func (l *List) shortOutput(data []string) {
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

// long prints the template list in long format.
func (l *List) long() int {
	templates := l.Conf.Templates
	if len(templates) == 0 {
		return 0
	}
	w := tabwriter.NewWriter(l.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	sort.Sort(config.SortByName(templates))

	var colorHeader ansicodes.Codes
	var colorStandard ansicodes.Codes
	var colorTemplate ansicodes.Codes
	if isatty.IsTerminal(os.Stdout.Fd()) {
		colorHeader = ansicodes.Faint
		colorStandard = ansicodes.Reset
		colorTemplate = ansicodes.GreenText
	}

	fmt.Fprintf(w, "%sHandle\t%sTemplate%s\tLocation%s\n", colorHeader, colorHeader, colorHeader, colorStandard)
	for _, stub := range templates {
		switch {
		case stub.Template == "" && stub.Origin == "":
			fmt.Fprintf(w, "%s%s\t%s%s\t%s%s\n", colorStandard, stub.Handle, colorStandard, colorStandard, stub.URL, colorStandard)
		case stub.Origin == "":
			fmt.Fprintf(w, "%s%s\t%s%s%s\t%s%s\n", colorStandard, stub.Handle, colorTemplate, stub.Template, colorStandard, stub.URL, colorStandard)
		default:
			fmt.Fprintf(w, "%s%s\t%s%s%s/%s\t%s%s\n", colorStandard, stub.Handle, colorTemplate, stub.Template, colorStandard, stub.Origin, stub.URL, colorStandard)
		}
	}
	w.Flush()
	return 0
}
