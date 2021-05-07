package formatter

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/kick-project/kick/internal/resources/ansicodes"
	"github.com/kick-project/kick/internal/services/search/entry"
	"github.com/mattn/go-isatty"
	term "github.com/wayneashleyberry/terminal-dimensions"
)

// Format Format is a interface to format search entries.
type Format interface {
	// Writer Writer is abstract type that takes a search.Entry then writes to io.Writer
	Writer(io.Writer, <-chan *entry.Entry)
}

//type Format func(<-chan *entry.Entry, io.Writer)

// Standard Standard is the default formatter for search.Entry
type Standard struct {
	// Disable ANSI Escape codes
	NoANSICodes bool
	noTTYFlag   *bool
}

func (s *Standard) noTTY() bool {
	if s.noTTYFlag != nil {
		return *s.noTTYFlag
	}
	noTTY := true
	if isatty.IsTerminal(os.Stdout.Fd()) {
		noTTY = false
	}
	s.noTTYFlag = &noTTY
	return *s.noTTYFlag
}

// Writer format suitable for standard output
func (s *Standard) Writer(w io.Writer, ch <-chan *entry.Entry) {
	min := uint(20)
	repeat := min
	y, _ := term.Height()
	if y > repeat {
		repeat = y
	}
	tabwr := tabwriter.NewWriter(w, 0, 0, 10, ' ', 0)
	i := uint(0)
	for e := range ch {
		// Header
		if i == 0 {
			if s.NoANSICodes || s.noTTY() {
				fmt.Fprint(tabwr, "Template\tLocation\n")
			} else {
				fmt.Fprintf(tabwr, "%vTemplate%v\tLocation%v\n", ansicodes.Faint, ansicodes.Faint, ansicodes.None)
			}
		}
		i++

		// Body
		if s.NoANSICodes || s.noTTY() {
			fmt.Fprintf(tabwr, "%s/%s\t%s\n", e.Name, e.RepoName, e.URL)
		} else {
			fmt.Fprintf(tabwr, "%v%s%v/%s\t%s%v\n", ansicodes.GreenText, e.Name, ansicodes.None, e.RepoName, e.URL, ansicodes.None)
		}

		// Re-print header
		if i == repeat {
			i = 0
		}
	}
	tabwr.Flush()
}
