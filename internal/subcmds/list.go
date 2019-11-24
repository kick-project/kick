package subcmds

import (
	"fmt"
	"github.com/crosseyed/prjstart/internal"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
	"os"
	"sort"
	"text/tabwriter"
)

// List starts the list sub command
func List(args []string) int {
	opts := internal.GetOptList(args)
	switch {
	case opts.Vars:
		return VariablesLongOutput()
	case opts.Remote:
		return ListRemote()
	case opts.All:
		ret1 := ListLocal(opts)
		ret2 := ListRemote()
		if ret1 > ret2 {
			return ret1
		}
		return ret2
	}
	return ListLocal(opts)
}

// ListLocal lists local templates
func ListLocal(opts *internal.OptList) int {
	switch {
	case opts.Url:
		VerboseOutput()
	default:
		ShortOutput()
	}
	return 0
}

func ShortOutput() {
	templates := internal.Config.Templates
	sort.Sort(internal.SortByName(templates))
	data := []string{}
	for _, t := range templates {
		name := t.Name
		data = append(data, name)
	}
	shortOutput(data)
}

func shortOutput(data []string) {
	termwidth, _ := terminal.Width()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	outstring := ""
	for _, name := range data {
		l := uint(0)
		if len(outstring) != 0 {
			l = uint(len(outstring))
			l += uint(len(name))
			l += 1
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

func VerboseOutput() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	templates := internal.Config.Templates
	sort.Sort(internal.SortByName(templates))
	for _, stub := range templates {
		fmt.Fprintf(w, "%s\t%s\n", stub.Name, stub.URL)
	}
	w.Flush()
}

func VariablesLongOutput() int {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	prjvars := internal.SetVars(&internal.OptStart{})
	for _, item := range prjvars.GetDescriptions() {
		fmt.Fprintf(w, ".Project.%s\t'%s'\n", item[0], item[1])
		w.Flush()
	}
	return 0
}

// ListRemote lists remote templates
func ListRemote() int {
	fmt.Println("Remote: ")
	for _, uri := range internal.Config.SetURLs {
		f := internal.NewFetcher(internal.Config)
		path := f.GetSet(uri)
		if path == "" {
			fmt.Printf("Cloud not fetch: %s\n", uri)
			continue
		}
		set := internal.LoadSet(path, "")
		for _, tmpl := range set.Templates {
			fmt.Printf("URI: %s\n", uri)
			fmt.Printf("  Name: %s\n", tmpl.Name)
			fmt.Printf("  URL: %s\n", tmpl.URL)
			fmt.Printf("  Desc: %s\n", tmpl.Desc)
		}
	}
	return 0
}

// ListAll Lists local and remote templates
func ListAll() int {
	return 0
}
