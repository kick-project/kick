package install

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/file"
	"github.com/kick-project/kick/internal/resources/gitclient"
	"github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/resources/parse"
	"github.com/kick-project/kick/internal/resources/sync"
	"gorm.io/gorm"
)

// Install manage installation of templates
type Install struct {
	ConfigFile *config.File           `validate:"required"`
	ORM        *gorm.DB               `validate:"required"`
	Log        logger.OutputIface     `validate:"required"`
	Err        *errs.Handler          `validate:"required"`
	Plumb      plumbing.PlumbingIface `validate:"required"`
	Stderr     io.Writer              `validate:"required"`
	Stdin      io.Reader              `validate:"required"`
	Stdout     io.Writer              `validate:"required"`
	Sync       sync.SyncIface         `validate:"required"`
}

var selectWithOrigin = `
SELECT
	template.name AS templateName,
	template.url AS templateURL,
	repo.name AS origin,
	template.desc AS desc
FROM template LEFT JOIN repo_template ON (template.id = repo_template.template_id)
LEFT JOIN repo ON (repo_template.repo_id = repo.id)
WHERE template.name = ? AND repo.name = ?
`

var selectWithoutOrigin = `
SELECT
	template.name AS templateName,
	template.url AS templateURL,
	repo.name AS origin,
	template.desc AS desc
FROM template LEFT JOIN repo_template ON (template.id = repo_template.template_id)
LEFT JOIN repo ON (repo_template.repo_id = repo.id)
WHERE template.name = ?
`

// Install install template
func (i *Install) Install(handle, template string) (ret int) {
	i.Log.Debugf("Install(%s, %s)", handle, template)

	// Check if handle is in use
	ret = i.checkInUse(handle)
	if ret != 0 {
		i.Log.Printf("handle %s is already in use\n", handle)
		return
	}

	// Install from a template name
	found := i.processTemplate(handle, template)
	if found {
		return
	}

	// Install from a URL
	found = i.processLocation(handle, template)
	if !found {
		i.Log.Printf("invalid template or url %s\n", template)
		ret = 255
	}
	return
}

func (i *Install) processLocation(handle, location string) (found bool) {
	i.Log.Debugf("processLocation(%s, %s)", handle, location)

	p, err := filepath.Abs(file.ExpandPath(location))
	i.Err.Panic(err)
	// Check if its a path on the local file system
	if info, err := os.Stat(p); err == nil && info.IsDir() {
		t := config.Template{
			URL: p,
		}
		i.createEntry(handle, t)
		return true
	}

	urlx, err := parse.Parse(location)
	if err != nil {
		return false
	}
	t := config.Template{
		URL:  urlx.URL,
		Desc: "Direct installation",
	}
	i.createEntry(handle, t)
	return true
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
	row := i.ORM.Raw(`SELECT count(*) AS count FROM installed WHERE ?`, handle).Row()
	err := row.Scan(&count)
	errs.Panic(err)

	if count > 0 {
		return 255
	}

	return 0
}

// templateMatches searches for template matches in the database and
// returns them as entries. If stop is returned as non 0 then the caller should
// exit the program execution with the value of stop.
func (i *Install) templateMatches(name, origin string) (entries []config.Template) {
	var rows *sql.Rows
	entries = []config.Template{}
	if origin == "" {
		r, err := i.ORM.Raw(selectWithoutOrigin, name).Rows()
		errs.Fatal(err)
		rows = r
	} else {
		r, err := i.ORM.Raw(selectWithOrigin, name, origin).Rows()
		errs.Fatal(err)
		rows = r
	}
	for rows.Next() {
		var (
			template sql.NullString
			URL      sql.NullString
			origin   sql.NullString
			desc     sql.NullString
		)
		err := rows.Scan(&template, &URL, &origin, &desc)
		errs.Panic(err)

		entry := config.Template{
			Template: template.String,
			URL:      URL.String,
			Origin:   origin.String,
			Desc:     desc.String,
		}

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

	var match []string
	selected := 0
	re := regexp.MustCompile(`^(\d+)\n$`)
	for {
		fmt.Fprintf(i.Stdout, "  Select an entry between 1-%d: ", l)
		reader := bufio.NewReader(i.Stdin)
		text, err := reader.ReadString('\n')
		i.Err.Panic(err)

		match = re.FindStringSubmatch(text)
		if len(match) == 0 {
			fmt.Fprint(i.Stdout, "\nInvalid entry\n\n")
		} else {
			selected, err = strconv.Atoi(match[1])
			i.Err.Panic(err)
		}

		if selected < 1 || selected > l {
			fmt.Fprint(i.Stdout, "\nInvalid entry\n\n")
		} else {
			break
		}
	}
	i.createEntry(handle, entries[selected-1])
}

// createEntry creates a entry
func (i *Install) createEntry(handle string, entry config.Template) {
	_, err := i.getRepo(entry.URL)
	i.Err.FatalF("Error installing %s: %v\n", entry.Handle, err)

	entry.Handle = handle
	err = i.ConfigFile.AppendTemplate(entry)
	i.Err.Fatal(err)
	err = i.ConfigFile.SaveTemplates()
	i.Err.Fatal(err)
	i.Sync.Files()
	switch {
	case entry.Template == "":
		i.Log.Printf("installed handle:%s -> location:%s\n", entry.Handle, entry.URL)
	case entry.Origin == "":
		i.Log.Printf("installed handle:%s template:%s -> location:%s\n", entry.Handle, entry.Template, entry.URL)
	default:
		i.Log.Printf("installed handle:%s template:%s/%s -> location:%s\n", entry.Handle, entry.Template, entry.Origin, entry.URL)
	}
}

// getRepo get version control system repository or set a location to a template.
// returns the local path location.
func (i *Install) getRepo(url string) (path string, err error) {
	path, err = gitclient.Get(url, i.Plumb)
	return
}
