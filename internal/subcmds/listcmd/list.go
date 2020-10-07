package listcmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/settings"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/docopt/docopt-go"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

var usageDoc = `List templates

Usage:
    prjstart list [--url]
    prjstart list [--vars]

Options:
    -h --help     Print help.
    -u --url      Print URL.
    -v --vars     Show Variables.
`

type OptList struct {
	List   bool `docopt:"list"`
	Local  bool `docopt:"--local"`
	Remote bool `docopt:"--remote"`
	All    bool `docopt:"--all"`
	URL    bool `docopt:"--url"`
	Vars   bool `docopt:"--vars"`
}

// GetOptStart parse start options from document text
func GetOptList(args []string) *OptList {
	opts, err := docopt.ParseArgs(usageDoc, args, "")
	errutils.Epanicf("Can not parse usage doc: %s", err) // nolint
	o := new(OptList)
	err = opts.Bind(o)
	errutils.Epanicf("Can not bind to structure: %s", err) // nolint
	return o
}

type listCmd struct {
	conf *config.File
}

// List starts the list sub command
func List(args []string, s *settings.Settings) int {
	opts := GetOptList(args)
	lc := listCmd{
		conf: s.ConfigFile(),
	}
	switch {
	case opts.Vars:
		return lc.VariablesLongOutput()
	case opts.Remote:
		return lc.ListRemote()
	case opts.All:
		ret1 := lc.ListLocal(opts)
		ret2 := lc.ListRemote()
		if ret1 > ret2 {
			return ret1
		}
		return ret2
	}
	return lc.ListLocal(opts)
}

// ListLocal lists local templates
func (lc *listCmd) ListLocal(opts *OptList) int {
	switch {
	case opts.URL:
		lc.VerboseOutput()
	default:
		lc.ShortOutput()
	}
	return 0
}

func (lc *listCmd) ShortOutput() {
	templates := lc.conf.TemplateURLs
	sort.Sort(config.SortByName(templates))
	data := []string{}
	for _, t := range templates {
		name := t.Name
		data = append(data, name)
	}
	lc.shortOutput(data)
}

func (lc *listCmd) shortOutput(data []string) {
	termwidth, _ := terminal.Width()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
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

func (lc *listCmd) VerboseOutput() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	templates := lc.conf.TemplateURLs
	sort.Sort(config.SortByName(templates))
	for _, stub := range templates {
		fmt.Fprintf(w, "%s\t%s\n", stub.Name, stub.URL)
	}
	w.Flush()
}

func (lc *listCmd) VariablesLongOutput() int {
	// w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	// prjvars := internal.SetVars(&start.OptStart{})
	// for _, item := range prjvars.GetDescriptions() {
	// 	fmt.Fprintf(w, ".Project.%s\t'%s'\n", item[0], item[1])
	// 	w.Flush()
	// }
	return 0
}

// ListRemote lists remote templates
func (lc *listCmd) ListRemote() int {
	// TODO
	return 0
}

// ListAll Lists local and remote templates
func (lc *listCmd) ListAll() int {
	return 0
}
