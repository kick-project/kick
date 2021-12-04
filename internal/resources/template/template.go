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

	"github.com/go-playground/validator"
	"github.com/kick-project/kick/internal/resources/checkvars"
	"github.com/kick-project/kick/internal/resources/client"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/exit"
	"github.com/kick-project/kick/internal/resources/file"
	"github.com/kick-project/kick/internal/resources/logger"
	"github.com/kick-project/kick/internal/resources/marshal"
	"github.com/kick-project/kick/internal/resources/template/renderer"
	"github.com/kick-project/kick/internal/resources/template/variables"
)

const (
	// MLnone file does not have a mode line.
	MLnone = iota
	// MLrender file mode line instruction to render.
	MLrender
	// MLnorender file mode line instruction not to render.
	MLnorender
)

// Template the template itself
//go:generate ifacemaker -f template.go -s Template -p template -i TemplateIface -o template_interfaces.go -c "AUTO GENERATED. DO NOT EDIT."
type Template struct {
	checkvars      *checkvars.Check
	client         *client.Client
	config         *config.File
	errs           *errs.Handler
	exit           *exit.Handler
	log            logger.OutputIface
	modeLineLen    uint8
	nounset        bool
	noempty        bool
	renderCurrent  string
	renderersAvail map[string]renderer.Renderer
	stderr         io.Writer
	stdout         io.Writer
	templateDir    string
	vars           *variables.Variables
	builddir       string
	dest           string
	localpath      string
	src            string
}

// Options options to constructor
type Options struct {
	Checkvars      *checkvars.Check             // Not required
	Client         *client.Client               `validate:"required"`
	Config         *config.File                 `validate:"required"`
	Errs           *errs.Handler                `validate:"required"`
	Exit           *exit.Handler                `validate:"required"`
	Log            logger.OutputIface           `validate:"required"`
	NoUnset        bool                         // No unset variables
	NoEmpty        bool                         // No empty variables
	RenderCurrent  string                       `validate:"required"`
	RenderersAvail map[string]renderer.Renderer `validate:"required"`
	Stderr         io.Writer                    `validate:"required"`
	Stdout         io.Writer                    `validate:"required"`
	TemplateDir    string                       `validate:"required"`
	Variables      *variables.Variables         // Not required
	ModeLineLen    uint8
}

// New *Template constructor
func New(opts *Options) *Template {
	var (
		modeLineLen uint8
		err         error
	)
	err = validator.New().Struct(opts)
	if err != nil {
		panic(err)
	}
	if opts.ModeLineLen == 0 {
		modeLineLen = 5
	} else {
		modeLineLen = opts.ModeLineLen
	}
	return &Template{
		checkvars:      opts.Checkvars,
		client:         opts.Client,
		config:         opts.Config,
		errs:           opts.Errs,
		exit:           opts.Exit,
		log:            opts.Log,
		modeLineLen:    modeLineLen,
		nounset:        opts.NoUnset,
		noempty:        opts.NoEmpty,
		renderCurrent:  opts.RenderCurrent,
		renderersAvail: opts.RenderersAvail,
		stderr:         opts.Stderr,
		stdout:         opts.Stdout,
		templateDir:    opts.TemplateDir,
		vars:           opts.Variables,
	}
}

func (t *Template) chkvars(fp string) {
	if _, err := os.Stat(fp); err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			t.errs.FatalF(`error opening %s: %w`, fp, err)
		}
	}

	f, err := os.Open(fp)
	t.errs.FatalF(`error opening %s: %w`, fp, err)
	ok, err := t.checkvars.Check(f)
	t.errs.FatalF(`error checking vars: %w`, err)
	if !ok {
		t.exit.Exit(-1)
	}
}

// SetRender set rendering engine
func (t *Template) SetRender(renderer string) {
	if renderer == "" {
		t.log.Error("No renderer provided\n")
		exit.Exit(255)
	}

	if _, ok := t.renderersAvail[t.renderCurrent]; !ok {
		t.log.Errorf("No such renderer %s. Valid options are...\n", t.renderCurrent)
		for r := range t.renderersAvail {
			fmt.Fprintln(t.stdout, r)
		}
		t.exit.Exit(255)
	}
	t.renderCurrent = renderer
}

func (t *Template) SetVars(vars *variables.Variables) {
	t.vars = vars
}

func (t *Template) renderer() renderer.Renderer {
	if t.renderCurrent == "" {
		panic("no render")
	}

	render, ok := t.renderersAvail[t.renderCurrent]
	if ok {
		return render
	}
	t.log.Printf("No such renderer %s\n", t.renderCurrent)
	t.exit.Exit(255)
	return nil
}

func (t *Template) buildDir(id string) {
	d, err := ioutil.TempDir(os.Getenv("TEMP"), fmt.Sprintf("kick-%s-", id))
	t.errs.PanicF("build error: %v", err)
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
	var tmpl config.Template
	for _, tconf := range t.config.Templates {
		if tconf.Handle == name {
			tmpl = tconf
			break
		}
	}
	p, err := t.client.GetTemplate(tmpl.URL, "")
	t.errs.FatalF(`handle "%s" not found: %v`, name, err)
	localpath := p.Path()

	// Check for missing variables
	y := filepath.Join(localpath, `.kick.yml`)
	t.chkvars(y)

	// Set renderer from conf
	confPath := filepath.Join(localpath, ".kick.yml")
	t.loadTempateConf(confPath)

	stat, err := os.Stat(localpath)
	t.errs.FatalF(`error: %w`, err)

	if !stat.IsDir() {
		fmt.Fprintf(os.Stderr, `%s is not a directory`, localpath)
		t.exit.Exit(-1)
	}

	t.src = name
	t.localpath = localpath
}

// SetDest sets the destination path
func (t *Template) SetDest(dest string) {
	t.dest = dest
}

// Run generates the target directory structure
func (t *Template) Run() int {
	path := t.localpath
	base := t.localpath
	skipRegex, err := regexp.Compile(fmt.Sprintf(`^%s/.git(?:/|$)`, base))
	t.errs.PanicF("build error: %v", err)
	t.checkDstExists()
	errWalk := filepath.Walk(path, func(srcPath string, info os.FileInfo, err error) error {
		if skipRegex.MatchString(srcPath) {
			return nil
		}
		relative := strings.Replace(srcPath, base, "", 1)
		relative = t.renderDir(relative)
		dstPath := filepath.Join(t.builddir, relative)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Warning: skipping file %s: %s", srcPath, err.Error())
			return nil
		}

		pair := filePair{
			errs:      t.errs,
			srcInfo:   info,
			srcPath:   srcPath,
			dstPath:   dstPath,
			variables: t.vars,
			mlen:      t.modeLineLen,
			renderer:  t.renderer(),
		}
		err = pair.route()
		t.errs.PanicF("build error: %v", err)

		return nil
	})
	if errWalk != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Abort creating project: %s", errWalk.Error())
		return 255
	}

	err = file.MoveAll(t.builddir, t.dest)
	t.errs.PanicF("build error: %v", err)
	t.log.Printf(`created project handle:%s -> project:%s`, t.src, t.dest)
	return 0
}

func (t *Template) checkDstExists() {
	stat, err := os.Stat(t.dest)
	if !os.IsNotExist(err) {
		t.errs.PanicF("build error: %v", err)
	}
	if stat != nil {
		t.log.Printf("path '%s' exists. aborting.\n", t.dest) // nolint
		t.exit.Exit(255)
	}
}

// renderDir scans directory names for template markers and renders the directory path as a template
func (t *Template) renderDir(path string) string {
	regex := t.renderer().RenderDirRegexp()
	if !regex.MatchString(path) {
		return path
	}
	path, err := t.renderer().Text2String(path, t.vars, true, true)
	t.errs.FatalF("can not substitute path string \"%s\": %v", path, err)
	return path
}

// loadTemplateConf loads template configuration
func (t *Template) loadTempateConf(path string) {
	t.log.Debugf("loading template conf %s\n", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	} else if err != nil {
		fmt.Fprintf(t.stderr, "Can not load file %s: %v\n", path, err)
	}
	c := &templateConf{}
	if !strings.HasSuffix(path, ".yml") {
		fmt.Fprintf(t.stderr, "Invalid file %s\n", path)
		t.exit.Exit(255)
	}
	t.log.Debugf("Unmarshal file %s\n", path)
	err := marshal.FromFile(c, path)
	if err != nil {
		fmt.Fprintf(t.stderr, "Can not unmarshal file %s: %s\n", path, err.Error())
		t.exit.Exit(255)
	}
	t.log.Debugf("%#v", c)

	if c.Renderer != "" {
		t.SetRender(c.Renderer)
	}
}

//
// Source Destination pair
//
type filePair struct {
	dstPath   string // Destination path
	mlen      uint8  // Mode line length
	errs      *errs.Handler
	mu        sync.Mutex
	nounset   bool
	noempty   bool
	renderer  renderer.Renderer
	srcInfo   os.FileInfo
	srcPath   string // Source path
	variables *variables.Variables
}

func (fp *filePair) route() error {
	action, lnum := fp.hasModeLine()
	switch {
	case fp.srcInfo.IsDir():
		return fp.mkdir()
	case fp.skipFile():
		return nil
	case lnum > 0 && action == MLrender:
		err := fp.render(lnum)
		fp.errs.Panic(err)
	case lnum > 0 && action == MLnorender:
		return nil
	case fp.srcInfo.Mode().IsRegular():
		return fp.copy()
	default:
		msg := fmt.Sprintf("error FILENOTREGULAR: %s\n", fp.dstPath)
		return errors.New(msg)
	}
	return nil
}

// skipFile determines known files to skip
func (fp *filePair) skipFile() bool {
	rvalue := false
	switch {
	case strings.HasSuffix(fp.srcPath, ".kick.yml"):
		rvalue = true
	}
	return rvalue
}

func (fp *filePair) mkdir() error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	if _, err := os.Stat(fp.dstPath); os.IsNotExist(err) {
		err = os.Mkdir(fp.dstPath, 0755)
		fp.errs.PanicF("build error: %v", err)
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
	fp.errs.PanicF("Can not open '%s': %s", fp.srcPath, err) // nolint
	defer inF.Close()                                        // nolint

	tmpdir := os.Getenv("TMPDIR")
	outF, err := ioutil.TempFile(tmpdir, "kick-")
	fp.errs.PanicF("Can not create tempfile: %s", err) // nolint
	defer outF.Close()                                 // nolint

	var cnt uint8
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
		fp.errs.PanicF("Error writing to file '%s': %s", outF.Name(), err) // nolint
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

	err := fp.renderer.File2File(tempPath, fp.dstPath, fp.variables, fp.nounset, fp.noempty)
	fp.errs.Panic(err)
	return nil
}

// hasModeLine scans the first mlen lines for a modeline
// returns MLnone, MLrender depending on defined action
func (fp *filePair) hasModeLine() (action int, lnum uint8) {
	action = MLnone
	len := fp.mlen
	mlactions := hasML{}.Init()
	source, err := os.Open(fp.srcPath)
	fp.errs.FatalF("Can not open file %s: %v", fp.srcPath, err)

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

//
// hasML
//

type hasML []hasMLAction

func (ml hasML) Init() hasML {
	ml = append(ml, regexCompile(`kick:render\W?`, MLrender))
	ml = append(ml, regexCompile(`kick:ignore\W?`, MLnorender))
	return ml
}

//
// templateConf
//

type templateConf struct {
	Renderer string `yaml:"renderer"`
}

func regexCompile(rex string, action int) hasMLAction {
	regex, err := regexp.Compile(rex)
	errs.PanicF("Error compiling regex: %s", err) // nolint
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
