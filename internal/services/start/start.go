package start

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/kick-project/kick/internal/resources/check"
	"github.com/kick-project/kick/internal/resources/checkvars"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/sync"
	"github.com/kick-project/kick/internal/resources/template"
	"github.com/kick-project/kick/internal/resources/template/variables"
	"github.com/olekukonko/tablewriter"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

// Start manage listing of installed templates
type Start struct {
	check     *check.Check
	checkvars *checkvars.Check
	conf      *config.File
	exit      exit.HandlerIface
	stderr    io.Writer
	stdout    io.Writer
	sync      sync.SyncIface
	tmpl      template.TemplateIface
}

// Options contructor options
type Options struct {
	Check     *check.Check           `validate:"required"`
	CheckVars *checkvars.Check       `validate:"required"`
	Conf      *config.File           `validate:"required"`
	Exit      exit.HandlerIface      `validate:"required"`
	Stderr    io.Writer              `validate:"required"`
	Stdout    io.Writer              `validate:"required"`
	Sync      sync.SyncIface         `validate:"required"`
	Template  template.TemplateIface `validate:"required"`
}

// New consructor
func New(opts Options) *Start {
	s := &Start{
		check:     opts.Check,
		checkvars: opts.CheckVars,
		conf:      opts.Conf,
		exit:      opts.Exit,
		stderr:    opts.Stderr,
		stdout:    opts.Stdout,
		sync:      opts.Sync,
		tmpl:      opts.Template,
	}
	return s
}

// Start start command
func (s Start) Start(projectname, template, path string) {
	if err := s.check.Init(); err != nil {
		fmt.Fprintf(s.stderr, "%s\n", err.Error())
		s.exit.Exit(255)
	}

	// Sync DB table "installed" with configuration file
	s.sync.Files()

	// Set varaibles
	vars := variables.New()
	vars.ProjectVariable("NAME", projectname)
	s.tmpl.SetVars(vars)

	// Set project name
	s.tmpl.SetSrcDest(template, path)
	_ = s.tmpl.Run()
}

// List lists the output
func (s *Start) List(long bool) {
	if long {
		s.longFmt()
	} else {
		s.shortFmt()
	}
}

func (s *Start) shortFmt() {
	templates := s.sort(s.conf.Templates)
	sort.Sort(config.SortByName(templates))
	data := []string{}
	for _, t := range templates {
		name := t.Handle
		data = append(data, name)
	}
	termwidth, _ := terminal.Width()
	w := tabwriter.NewWriter(s.stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
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

func (s *Start) longFmt() {
	var (
		header []string
		table  [][]string
	)
	header = []string{"Handle", "Template", "Description", "Location"}
	for _, row := range s.sort(s.conf.Templates) {
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
		table = append(table, []string{row.Handle, templateName, desc, row.URL})
	}
	writer := tablewriter.NewWriter(s.stdout)
	writer.SetAlignment(tablewriter.ALIGN_LEFT)
	writer.SetHeader(header)
	for _, v := range table {
		writer.Append(v)
	}
	writer.Render()
}

func (s *Start) sort(in []config.Template) (out []config.Template) {
	sort.Sort(config.SortByName(in))
	out = append(out, in...)
	return
}
