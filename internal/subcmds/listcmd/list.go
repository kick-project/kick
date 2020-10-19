package listcmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/jinzhu/copier"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/settings"
	"github.com/kick-project/kick/internal/settings/isync"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/options"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

var usageDoc = `List handles/templates

Usage:
    prjstart list [-l]

Options:
	-h --help     Print help.
    -l            Print Long output.
`

type OptList struct {
	List bool `docopt:"list"`
	Long bool `docopt:"-l"`
}

type listCmd struct {
	conf *config.File
}

// List starts the list sub command
func List(args []string, s *settings.Settings) int {
	opts := &OptList{}
	options.Bind(usageDoc, args, opts)
	lc := listCmd{
		conf: s.ConfigFile(),
	}

	// Sync DB table "installed" with configuration file
	synchro := &sync.Sync{}
	err := copier.Copy(synchro, isync.Inject(s))
	errutils.Epanic(err)
	synchro.Templates()
	if opts.Long {
		return lc.LongOutput()
	}
	return lc.ShortOutput()
}

func (lc *listCmd) ShortOutput() int {
	templates := lc.conf.Templates
	sort.Sort(config.SortByName(templates))
	data := []string{}
	for _, t := range templates {
		name := t.Handle
		data = append(data, name)
	}
	lc.shortOutput(data)
	return 0
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

func (lc *listCmd) LongOutput() int {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	templates := lc.conf.Templates
	sort.Sort(config.SortByName(templates))
	for _, stub := range templates {
		fmt.Fprintf(w, "%s\t%s\n", stub.Handle, stub.URL)
	}
	w.Flush()
	return 0
}
