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

	"github.com/apex/log"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/gitclient"
	plumb "github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/resources/template/renderer"
	"github.com/kick-project/kick/internal/resources/template/variables"
	"github.com/kick-project/kick/internal/utils"
	"github.com/kick-project/kick/internal/utils/errutils"
	"github.com/kick-project/kick/internal/utils/marshal"
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
type Template struct {
	Config         *config.File
	Log            *log.Logger
	ModeLineLen    uint8
	RenderCurrent  string
	RenderersAvail map[string]renderer.Renderer
	Stderr         io.Writer
	Stdout         io.Writer
	TemplateDir    string
	Variables      *variables.Variables
	builddir       string
	dest           string
	localpath      string
	src            string
}

// SetRender set rendering engine
func (t *Template) SetRender(renderer string) {
	if renderer == "" {
		t.Log.Error("No renderer provided\n")
		utils.Exit(255)
	}

	if _, ok := t.RenderersAvail[t.RenderCurrent]; !ok {
		t.Log.Errorf("No such renderer %s. Valid options are...\n", t.RenderCurrent)
		for r := range t.RenderersAvail {
			fmt.Println(r)
		}
		utils.Exit(255)
	}
	t.RenderCurrent = renderer
}

func (t *Template) renderer() renderer.Renderer {
	if t.RenderCurrent == "" {
		panic("no render")
	}

	render, ok := t.RenderersAvail[t.RenderCurrent]
	if ok {
		return render
	} else {
		fmt.Fprintf(t.Stderr, "No such renderer %s\n", t.RenderCurrent)
		utils.Exit(255)
	}
	return nil
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
	var tmpl config.Template
	for _, tconf := range t.Config.Templates {
		if tconf.Handle == name {
			tmpl = tconf
			break
		}
	}
	g := plumb.New(t.TemplateDir)
	localpath, err := gitclient.Get(tmpl.URL, g)

	// Set renderer from conf
	confPath := filepath.Join(localpath, ".kicktemplate.yml")
	t.loadTempateConf(confPath)

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

// SetDest sets the destination path
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
		relative = t.renderDir(relative)
		dstPath := filepath.Join(t.builddir, relative)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Warning: skipping file %s: %s", srcPath, err.Error())
			return nil
		}

		pair := filePair{
			srcInfo:   info,
			srcPath:   srcPath,
			dstPath:   dstPath,
			variables: t.Variables,
			mlen:      t.ModeLineLen,
			renderer:  t.renderer(),
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
		fmt.Printf("Path '%s' exists. Aborting.\n", t.dest) // nolint
		utils.Exit(255)
	}
}

// renderDir scans directory names for template markers and renders the directory path as a template
func (t *Template) renderDir(path string) string {
	regex := t.renderer().RenderDirRegexp()
	if !regex.MatchString(path) {
		return path
	}
	path, err := t.renderer().Text2String(path, t.Variables, true, true)
	errutils.Efatalf("can not substitute path string \"%s\": %v", path, err)
	return path
}

// loadTemplateConf loads template configuration
func (t *Template) loadTempateConf(path string) {
	t.Log.Debugf("Loading template conf %s\n", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	} else if err != nil {
		fmt.Fprintf(t.Stderr, "Can not load file %s: %v\n", path, err)
	}
	c := &templateConf{}
	if !strings.HasSuffix(path, ".yml") {
		fmt.Fprintf(t.Stderr, "Invalid file %s\n", path)
		utils.Exit(255)
	}
	t.Log.Debugf("Unmarshal file %s\n", path)
	marshal.UnmarshalFile(c, path)
	t.Log.Debugf("%#v", c)

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
	mu        sync.Mutex
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
		fp.render(lnum)
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
	case strings.HasSuffix(fp.srcPath, ".kickglobal.yml"):
		rvalue = true
	case strings.HasSuffix(fp.srcPath, ".kickmaster.yml"):
		rvalue = true
	case strings.HasSuffix(fp.srcPath, ".kicktemplate.yml"):
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

	fp.renderer.File2File(tempPath, fp.dstPath, fp.variables, false, false)
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

//
// hasML
//

type hasML []hasMLAction

func (ml hasML) Init() hasML {
	ml = append(ml, regexCompile(`prj:render\W?`, MLrender))
	ml = append(ml, regexCompile(`prj:ignore\W?`, MLnorender))
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
