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

	"github.com/kick-project/kick/internal/resources/client"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/file"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/resources/parse"
	"github.com/kick-project/kick/internal/resources/sync"
	"gorm.io/gorm"
)

// Install manage installation of templates
type Install struct {
	client     *client.Client
	ConfigFile *config.File
	orm        *gorm.DB
	log        logger.OutputIface
	exit       *exit.Handler
	err        *errs.Handler
	stderr     io.Writer
	stdin      io.Reader
	stdout     io.Writer
	sync       sync.SyncIface
}

// Options options to constructor
type Options struct {
	Client     *client.Client     `validate:"required"`
	ConfigFile *config.File       `validate:"required"`
	ORM        *gorm.DB           `validate:"required"`
	Log        logger.OutputIface `validate:"required"`
	Exit       *exit.Handler      `validate:"required"`
	Err        *errs.Handler      `validate:"required"`
	Stderr     io.Writer          `validate:"required"`
	Stdin      io.Reader          `validate:"required"`
	Stdout     io.Writer          `validate:"required"`
	Sync       sync.SyncIface     `validate:"required"`
}

// New constructor
func New(opts *Options) *Install {
	return &Install{
		client:     opts.Client,
		ConfigFile: opts.ConfigFile,
		orm:        opts.ORM,
		log:        opts.Log,
		exit:       opts.Exit,
		err:        opts.Err,
		stderr:     opts.Stderr,
		stdin:      opts.Stdin,
		stdout:     opts.Stdout,
		sync:       opts.Sync,
	}
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
	i.log.Debugf("Install(%s, %s)", handle, template)

	// Check if handle is in use
	ret = i.checkInUse(handle)
	if ret != 0 {
		i.log.Printf("handle %s is already in use\n", handle)
		return
	}

	// Install from a template name
	found := i.processTemplate(handle, template)
	if found {
		return 0
	}

	// Install from a URL
	found, err := i.processLocation(handle, template)
	if err != nil {
		return 255
	} else if !found {
		i.log.Printf("invalid template or url %s\n", template)
		ret = 255
	}
	return
}

func (i *Install) processLocation(handle, location string) (found bool, err error) {
	i.log.Debugf("processLocation(%s, %s)", handle, location)

	p, err := filepath.Abs(file.ExpandPath(location))
	if err != nil {
		return false, err
	}
	// Check if its a path on the local file system
	if info, err := os.Stat(p); err == nil && info.IsDir() {
		t := config.Template{
			URL: p,
		}
		err = i.createEntry(handle, t)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	urlx, err := parse.Parse(location)
	if err != nil {
		return false, err
	}
	t := config.Template{
		URL:  urlx.URL,
		Desc: "Direct installation",
	}
	err = i.createEntry(handle, t)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (i *Install) processTemplate(handle, template string) (processed bool) {
	i.log.Debugf("processTemplate(%s, %s)", handle, template)
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
		_ = i.createEntry(handle, entries[0])
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
	row := i.orm.Raw(`SELECT count(*) AS count FROM installed WHERE ?`, handle).Row()
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
		r, err := i.orm.Raw(selectWithoutOrigin, name).Rows()
		errs.Fatal(err)
		rows = r
	} else {
		r, err := i.orm.Raw(selectWithOrigin, name, origin).Rows()
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
	fmt.Fprint(i.stdout, "multiple matches\n", l)
	for x := 0; x < l; x++ {
		cur := entries[x]
		fmt.Fprintf(i.stdout, "  (%d): %s/%s %s\n", x+1, cur.Handle, cur.Origin, cur.URL)
	}
	fmt.Fprint(i.stdout, "\n  Please select an entry\n")

	var match []string
	selected := 0
	re := regexp.MustCompile(`^(\d+)\n$`)
	for {
		fmt.Fprintf(i.stdout, "  Select an entry between 1-%d: ", l)
		reader := bufio.NewReader(i.stdin)
		text, err := reader.ReadString('\n')
		i.err.Panic(err)

		match = re.FindStringSubmatch(text)
		if len(match) == 0 {
			fmt.Fprint(i.stdout, "\nInvalid entry\n\n")
		} else {
			selected, err = strconv.Atoi(match[1])
			i.err.Panic(err)
		}

		if selected < 1 || selected > l {
			fmt.Fprint(i.stdout, "\nInvalid entry\n\n")
		} else {
			break
		}
	}
	_ = i.createEntry(handle, entries[selected-1])
}

// createEntry creates a entry
func (i *Install) createEntry(handle string, entry config.Template) error {
	_, err := i.getRepo(entry.URL)
	if err != nil {
		return err
	}

	entry.Handle = handle
	err = i.ConfigFile.AppendTemplate(entry)
	if err != nil {
		return fmt.Errorf(`entry error: %w`, err)
	}
	err = i.ConfigFile.SaveTemplates()
	if err != nil {
		return fmt.Errorf(`entry error: %w`, err)
	}
	i.sync.Files()
	switch {
	case entry.Template == "":
		i.log.Printf("installed handle:%s -> location:%s\n", entry.Handle, entry.URL)
	case entry.Origin == "":
		i.log.Printf("installed handle:%s template:%s -> location:%s\n", entry.Handle, entry.Template, entry.URL)
	default:
		i.log.Printf("installed handle:%s template:%s/%s -> location:%s\n", entry.Handle, entry.Template, entry.Origin, entry.URL)
	}
	return nil
}

// getRepo get version control system repository or set a location to a template.
// returns the local path location.
func (i *Install) getRepo(url string) (string, error) {
	p, err := i.client.GetTemplate(url, "")
	if err != nil {
		return "", err
	}
	return p.Path(), err
}
