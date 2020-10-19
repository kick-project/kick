package install

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/apex/log"
	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/resources/gitclient"
	"github.com/crosseyed/prjstart/internal/resources/gitclient/plumbing"
	"github.com/crosseyed/prjstart/internal/resources/sync"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

// Install manage installation of templates
type Install struct {
	ConfigFile *config.File
	DB         *sql.DB
	Log        *log.Logger
	Plumb      *plumbing.Plumbing
	Stderr     io.Writer
	Stdin      io.Reader
	Stdout     io.Writer
	Sync       *sync.Sync
}

var selectWithOrigin = `
SELECT
	templates.name AS templateName,
	templates.url AS templateURL,
	master.name AS origin,
	templates.desc AS desc
FROM templates LEFT JOIN master ON (templates.masterid = master.id)
WHERE templates.name = ? AND master.name = ?
`

var selectWithoutOrigin = `
SELECT
	templates.name AS templateName,
	templates.url AS templateURL,
	master.name AS origin,
	templates.desc AS desc
FROM templates LEFT JOIN master ON (templates.masterid = master.id)
WHERE templates.name = ?
`

// Install install template
func (i *Install) Install(handle, template string) (ret int) {
	i.Log.Debugf("Install(%s, %s)", handle, template)

	// Check if handle is in use
	ret = i.checkInUse(handle)
	if ret != 0 {
		fmt.Fprintf(i.Stderr, "handle %s is already in use\n", handle)
		return
	}

	// Install from a template name
	found := i.processTemplate(handle, template)
	if found {
		return
	}

	// Install from a URL
	ret = i.processURL(handle, template)
	if ret != 0 {
		fmt.Fprintf(i.Stderr, "invalid template or url %s\n", template)
	}
	return
}

func (i *Install) processURL(handle, url string) (stop int) {
	i.Log.Debugf("processURL(%s, %s)", handle, url)
	urlx, err := utils.Parse(url)
	if err != nil {
		return 255
	}
	t := config.Template{
		URL:  urlx.URL,
		Desc: "Direct installation",
	}
	i.createEntry(handle, t)
	return
}

func (i *Install) processTemplate(handle, template string) (processed bool) {
	i.Log.Debugf("processTemplate(%s, %s)", handle, template)
	var (
		entries []config.Template
		full    string
		name    string
		origin  string
	)
	re := regexp.MustCompile(`^([a-z0-9]+)(?:/([a-z0-9]+))?$`)
	match := re.FindStringSubmatch(template)
	if len(match) == 0 {
		return
	}
	full = match[0]
	name = match[1]
	origin = match[2]
	if full == "" {
		return
	}

	// Add entry
	entries = i.templateMatches(name, origin)
	switch len(entries) {
	case 0:
		return false
	case 1:
		i.createEntry(handle, entries[0])
		return true
	default:
		i.promptEntry(handle, entries)
		return true
	}
}

// checkInUse Check if a handle is in use. If stop is non 0 then the caller
// should stop program execution.
func (i *Install) checkInUse(handle string) (stop int) {
	var (
		count int
	)
	row := i.DB.QueryRow(`SELECT count(*) AS count FROM installed WHERE ?`, handle)
	err := row.Scan(&count)
	errutils.Epanic(err)

	if count > 0 {
		return 255
	}

	return 0
}

// templateMatches searches for template matches in the database and
// returns them as entries. If stop is returned as non 0 then the caller should
// exit the program execution with the value of stop.
func (i *Install) templateMatches(name, origin string) (entries []config.Template) {
	rows := &sql.Rows{}
	entries = []config.Template{}
	if origin == "" {
		i.Log.Debugf(utils.SQL2fmt(selectWithoutOrigin), name)
		r, err := i.DB.Query(selectWithoutOrigin, name)
		errutils.Efatal(err)
		rows = r
	} else {
		r, err := i.DB.Query(selectWithOrigin, name, origin)
		errutils.Efatal(err)
		rows = r
	}
	for rows.Next() {
		var (
			entry config.Template
		)
		err := rows.Scan(&entry.Template, &entry.URL, &entry.Origin, &entry.Desc)
		errutils.Epanic(err)

		entries = append(entries, entry)
	}
	return entries
}

// promptEntry prompts for an entry
func (i *Install) promptEntry(handle string, entries []config.Template) {
	l := len(entries)
	fmt.Fprint(i.Stdout, "multiple matches\n", l)
	for x := 0; x < l; x++ {
		cur := entries[x]
		fmt.Fprintf(i.Stdout, "  (%d): %s/%s %s\n", x+1, cur.Handle, cur.Origin, cur.URL)
	}
	fmt.Fprint(i.Stdout, "\n  Please select an entry\n")

	match := []string{}
	selected := 0
	re := regexp.MustCompile(`^(\d+)\n$`)
	for {
		fmt.Fprintf(i.Stdout, "  Select an entry between 1-%d: ", l)
		reader := bufio.NewReader(i.Stdin)
		text, err := reader.ReadString('\n')

		match = re.FindStringSubmatch(text)
		if len(match) == 0 {
			fmt.Print("\nInvalid entry\n\n")
		} else {
			selected, err = strconv.Atoi(match[1])
			errutils.Epanic(err)
		}

		if selected < 1 || selected > l {
			fmt.Print("\nInvalid entry\n\n")
		} else {
			break
		}
	}
	i.createEntry(handle, entries[selected-1])
}

// createEntry creates a entry
func (i *Install) createEntry(handle string, entry config.Template) {
	_, err := i.getRepo(entry.URL)
	if err != nil {
		fmt.Fprintf(i.Stderr, "Error installing %s: %s\n", entry.Handle, err.Error())
		utils.Exit(255)
	}

	entry.Handle = handle
	err = i.ConfigFile.AppendTemplate(entry)
	if err != nil {
		fmt.Fprintf(i.Stderr, "%s\n", err.Error())
		utils.Exit(255)
	}
	err = i.ConfigFile.SaveTemplates()
	if err != nil {
		fmt.Fprintf(i.Stderr, "%s\n", err.Error())
		utils.Exit(255)
	}
	i.Sync.Templates()
	switch {
	case entry.Template == "":
		fmt.Fprintf(i.Stdout, "Installed %s -> %s\n", entry.Handle, entry.URL)
	case entry.Origin == "":
		fmt.Fprintf(i.Stdout, "Installed %s %s -> %s\n", entry.Handle, entry.Template, entry.URL)
	default:
		fmt.Fprintf(i.Stdout, "Installed %s %s/%s -> %s\n", entry.Handle, entry.Template, entry.Origin, entry.URL)
	}
}

// getRepo get version control system repository or set a location to a template.
// returns the local path location.
func (i *Install) getRepo(url string) (path string, err error) {
	path, err = gitclient.Get(url, i.Plumb)
	return
}
