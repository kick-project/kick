package template

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/resources/gitclient"
	plumb "github.com/crosseyed/prjstart/internal/resources/gitclient/plumbing"
	"github.com/crosseyed/prjstart/internal/services/template/variables"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

const (
	// MLnone file does not have a mode line.
	MLnone = iota
	// MLrender file mode line instruction to render.
	MLrender
	// MLnorender file mode line instruction not to render.
	MLnorender
)

type (
	// Template the template itself
	Template struct {
		localpath   string
		src         string
		dest        string
		builddir    string
		config      *config.File
		variables   *variables.Variables
		templateDir string
		mllen       uint8
	}

	// Options options to template
	Options struct {
		Config    *config.File
		Variables *variables.Variables

		// TemplateDir is the directory to store the downloaded templates.
		TemplateDir string

		// ModeLineLen is the number of lines to scan in a document to fetch the modeline.
		// If set to 0 then modeline defaults to 20 lines.
		ModeLineLen uint8
	}
)

// New constructs a Template. New will panic if any options are missing.
func New(opts Options) *Template {
	if opts.Config == nil {
		panic("opts.Config can not be nil")
	}
	if opts.TemplateDir == "" {
		panic("opts.TemplateDir can not be an empty string")
	}
	var mllen = uint8(20)
	if opts.ModeLineLen > 0 {
		mllen = opts.ModeLineLen
	}
	t := &Template{
		config:      opts.Config,
		templateDir: opts.TemplateDir,
		variables:   opts.Variables,
		mllen:       mllen,
	}
	return t
}

func (t *Template) buildDir(id string) {
	d, err := ioutil.TempDir(os.Getenv("TEMP"), fmt.Sprintf("prjstart-%s-", id))
	errutils.Epanicf("Build Error: %v", err)
	t.builddir = d
}

// SetSrcDest sets the source template and destination path where the project structure
// will reside.
func (t *Template) SetSrcDest(src, dest string) {
	t.SetSrc(src)
	t.SetDest(dest)
}

// SetSrc sets the source template "name". "name" is defined
// in *config.Config.TemplateURLs. *config.Config is provided as an Option to New.
func (t *Template) SetSrc(name string) {
	t.buildDir(name)
	var tmpl config.TemplateStub
	for _, tconf := range t.config.TemplateURLs {
		if tconf.Name == name {
			tmpl = tconf
			break
		}
	}
	g := plumb.New(t.templateDir)
	// TODO: DI
	localpath, err := gitclient.Get(tmpl.URL, g)
	errutils.Efatalf(`template "%s" not found: %v`, name, err)

	stat, err := os.Stat(localpath)
	errutils.Efatalf(`error: %w`, err)

	if !stat.IsDir() {
		fmt.Fprintf(os.Stderr, `%s is not a directory`, localpath)
		utils.Exit(-1)
	}

	t.src = name
	t.localpath = localpath
}

// SetDest sets the destnation path
func (t *Template) SetDest(dest string) {
	t.dest = dest
}

// Run generates the target directory structure
func (t *Template) Run() int {
	path := t.localpath
	base := t.localpath
	skipRegex, err := regexp.Compile(fmt.Sprintf(`^%s/.git(?:/|$)`, base))
	errutils.Epanicf("build error: %v", err)
	t.checkDstExists()
	errWalk := filepath.Walk(path, func(srcPath string, info os.FileInfo, err error) error {
		if skipRegex.MatchString(srcPath) {
			return nil
		}
		relative := strings.Replace(srcPath, base, "", 1)
		relative = renderDir(relative, t.variables)
		dstPath := filepath.Join(t.builddir, relative)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Warning: skipping file %s: %s", srcPath, err.Error())
			return nil
		}

		pair := filePair{
			srcInfo:   info,
			srcPath:   srcPath,
			dstPath:   dstPath,
			variables: t.variables,
			mlen:      t.mllen,
		}
		err = pair.route()
		errutils.Epanicf("Build Error: %v", err)

		return nil
	})
	if errWalk != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Abort creating project: %s", errWalk.Error())
		return 255
	}

	err = os.Rename(t.builddir, t.dest)
	errutils.Epanicf("Build Error: %v", err)
	return 0
}

func (t *Template) checkDstExists() {
	stat, err := os.Stat(t.dest)
	if !os.IsNotExist(err) {
		errutils.Epanicf("Build Error: %v", err)
	}
	if stat != nil {
		fmt.Printf("Path '%s' exists. Aborting.", t.dest) // nolint
		utils.Exit(255)
	}
}

//
// Source Destination pair
//
type filePair struct {
	srcInfo   os.FileInfo
	srcPath   string // Source path
	dstPath   string // Destination path
	mlen      uint8  // Mode line length
	variables *variables.Variables
	mu        sync.Mutex
}

func (fp *filePair) route() error {
	action, lnum := fp.hasModeLine()
	switch {
	case fp.srcInfo.IsDir():
		return fp.mkdir()
	case fp.skipFile():
		return nil
	case lnum > 0 && action == MLrender:
		fp.render(lnum)
	case lnum > 0 && action == MLnorender:
		return nil
	case fp.srcInfo.Mode().IsRegular():
		return fp.copy()
	default:
		msg := fmt.Sprintf("error FILENOTREGULAR: %s\n", fp.dstPath)
		fmt.Println(msg)
		return errors.New(msg)
	}
	return nil
}

// skipFile determines known files to skip
func (fp *filePair) skipFile() bool {
	rvalue := false
	switch {
	case strings.HasSuffix(fp.srcPath, ".prjglobal.yml"):
		rvalue = true
	case strings.HasSuffix(fp.srcPath, ".prjmaster.yml"):
		rvalue = true
	case strings.HasSuffix(fp.srcPath, ".prjtemplate.yml"):
		rvalue = true
	}
	return rvalue
}

func (fp *filePair) mkdir() error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	if _, err := os.Stat(fp.dstPath); os.IsNotExist(err) {
		err = os.Mkdir(fp.dstPath, 0755)
		errutils.Epanicf("Build Error: %v", err)
	}
	return nil
}

func (fp *filePair) copy() error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	sourceFileStat, err := os.Stat(fp.srcPath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", fp.srcPath)
	}

	source, err := os.Open(fp.srcPath)
	if err != nil {
		return err
	}
	defer source.Close() // nolint

	destination, err := os.Create(fp.dstPath)
	if err != nil {
		return err
	}
	defer destination.Close() // nolint
	_, err = io.Copy(destination, source)
	return err
}

func (fp *filePair) stripModeline(lnum uint8) string {
	inF, err := os.Open(fp.srcPath)
	errutils.Epanicf("Can not open '%s': %s", fp.srcPath, err) // nolint
	defer inF.Close()                                          // nolint

	tmpdir := os.Getenv("TMPDIR")
	outF, err := ioutil.TempFile(tmpdir, "prjstart-")
	errutils.Epanicf("Can not create tempfile: %s", err) // nolint
	defer outF.Close()                                   // nolint

	var cnt uint8
	cnt = 0
	scner := bufio.NewScanner(inF)
	scner.Split(scanLines)
	for scner.Scan() {
		b := scner.Bytes()
		// Strip modeline
		if cnt < lnum {
			cnt++
			if lnum == cnt {
				continue
			}
		}
		_, err := outF.Write(b)
		errutils.Epanicf("Error writing to file '%s': %s", outF.Name(), err) // nolint
	}
	return outF.Name()
}

func (fp *filePair) render(mline uint8) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	// Remove modeline
	tempPath := fp.stripModeline(mline)
	defer func() {
		os.Remove(tempPath)
	}()

	File2File(tempPath, fp.dstPath, fp.variables)
	return nil
}

// hasModeLine scans the first mlen lines for a modeline
// returns MLnone, MLrender depending on defined action
func (fp *filePair) hasModeLine() (action int, lnum uint8) {
	action = MLnone
	len := fp.mlen
	mlactions := hasML{}.Init()
	source, err := os.Open(fp.srcPath)
	errutils.Efatalf("Can not open file %s: %v", fp.srcPath, err)

	defer source.Close()
	scner := bufio.NewScanner(source)
LOOP:
	for scner.Scan() {
		lnum++
		line := scner.Bytes()
		for _, mlaction := range mlactions {
			hasMatch := mlaction.Regex.Match(line)
			if hasMatch {
				action = mlaction.Action
				break LOOP
			}
		}
		if lnum > len {
			lnum = 0
			break
		}
	}
	return action, lnum
}

type hasML []hasMLAction

func (ml hasML) Init() hasML {
	ml = append(ml, regexCompile(`prj:render\W?`, MLrender))
	ml = append(ml, regexCompile(`prj:ignore\W?`, MLnorender))
	return ml
}

func regexCompile(rex string, action int) hasMLAction {
	regex, err := regexp.Compile(rex)
	errutils.Epanicf("Error compiling regex: %s", err) // nolint
	ml := hasMLAction{
		Regex:  regex,
		Action: action,
	}
	return ml
}

// hasMLAction contains a regex and an action of MLignore, MLrender or MLnone
type hasMLAction struct {
	Regex  *regexp.Regexp
	Action int
}

// scanLines is a split function for a Scanner that returns each line of
// text, WITHOUT stripping any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r?\n`.
// The last non-empty line of input will be returned even if it has no
// newline.
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// renderDir scans directory names for template markers and renders the directory path as a template
func renderDir(path string, prjvars *variables.Variables) string {
	regex := regexp.MustCompile(`{{[^}}]+}}`)
	if !regex.MatchString(path) {
		return path
	}
	path = Txt2String(path, prjvars)
	return path
}
