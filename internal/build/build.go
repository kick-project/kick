package build

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

	"github.com/crosseyed/prjstart/internal/gitclient/tclient"
	"github.com/crosseyed/prjstart/internal/globals"
	"github.com/crosseyed/prjstart/internal/template"
	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
)

const (
	MLnone = iota
	MLrender
	MLignore
)

type Build struct {
	localpath string
	src       string
	dest      string
	builddir  string
}

func (s *Build) buildDir(id string) {
	d, err := ioutil.TempDir(os.Getenv("TEMP"), fmt.Sprintf("prjstart-%s-", id))
	errutils.Epanicf(err, "Build Error: %v", err)
	s.builddir = d
}

// SetSrc sets the source template and localpath
func (s *Build) SetSrc(src string) {
	s.buildDir(src)
	dwnldr := tclient.TClient{
		Config: globals.Config,
	}
	localpath := dwnldr.Get(src)
	if localpath == "" {
		fmt.Fprintf(os.Stderr, `template "%s" not found`, src) // nolint
		utils.Exit(-1)
	}

	stat, err := os.Stat(localpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, `error: %s`, err.Error())
		utils.Exit(-1)
	}
	if !stat.IsDir() {
		fmt.Fprintf(os.Stdout, `%s is not a directory`, localpath)
		utils.Exit(-1)
	}

	s.src = src
	s.localpath = localpath
}

func (s *Build) SetDest(dest string) {
	s.dest = dest
}

func (s *Build) Run() int {
	path := s.localpath
	base := s.localpath
	skipRegex, err := regexp.Compile(fmt.Sprintf(`^%s/.git(?:/|$)`, base))
	errutils.Epanicf(err, "Build Error: %v", err)
	s.checkDstExists()
	errWalk := filepath.Walk(path, func(srcPath string, info os.FileInfo, err error) error {
		if skipRegex.MatchString(srcPath) {
			return nil
		}
		relative := strings.Replace(srcPath, base, "", 1)
		relative = renderDir(relative, globals.Vars)
		dstPath := filepath.Join(s.builddir, relative)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Warning: skipping file %s: %s", srcPath, err.Error())
			return nil
		}

		pair := filePair{
			srcInfo: info,
			srcPath: srcPath,
			dstPath: dstPath,
			mlen:    20, // TODO - Hardcoded modeline
		}
		err = pair.route()
		errutils.Epanicf(err, "Build Error: %v", err)

		return nil
	})
	if errWalk != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Abort creating project: %s", errWalk.Error())
		return 255
	}

	err = os.Rename(s.builddir, s.dest)
	errutils.Epanicf(err, "Build Error: %v", err)
	return 0
}

func (s *Build) checkDstExists() {
	stat, err := os.Stat(s.dest)
	if !os.IsNotExist(err) {
		errutils.Epanicf(err, "Build Error: %v", err)
	}
	if stat != nil {
		fmt.Printf("Path '%s' exists. Aborting.", s.dest) // nolint
		utils.Exit(255)
	}
}

//
// Source Destination pair
//
type filePair struct {
	srcInfo os.FileInfo
	srcPath string // Source path
	dstPath string // Destination path
	mlen    uint8  // Mode line length
	mu      sync.Mutex
}

func (s *filePair) route() error {
	action, lnum := s.hasModeLine()
	switch {
	case s.srcInfo.IsDir():
		return s.mkdir()
	case s.skipFile():
		return nil
	case lnum > 0 && action == MLrender:
		s.render(lnum)
	case lnum > 0 && action == MLignore:
		return nil
	case s.srcInfo.Mode().IsRegular():
		return s.copy()
	default:
		msg := fmt.Sprintf("error FILENOTREGULAR: %s\n", s.dstPath)
		fmt.Println(msg)
		return errors.New(msg)
	}
	return nil
}

// skipFile determines known files to skip
func (s *filePair) skipFile() bool {
	rvalue := false
	switch {
	case strings.HasSuffix(s.srcPath, ".prjmaster.yml"):
		rvalue = true
	case strings.HasSuffix(s.srcPath, ".prjorg.yml"):
		rvalue = true
	case strings.HasSuffix(s.srcPath, ".prjmod.yml"):
		rvalue = true
	}
	return rvalue
}

func (s *filePair) mkdir() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := os.Stat(s.dstPath); os.IsNotExist(err) {
		err = os.Mkdir(s.dstPath, 0777)
		errutils.Epanicf(err, "Build Error: %v", err)
	}
	return nil
}

func (s *filePair) copy() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sourceFileStat, err := os.Stat(s.srcPath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", s.srcPath)
	}

	source, err := os.Open(s.srcPath)
	if err != nil {
		return err
	}
	defer source.Close() // nolint

	destination, err := os.Create(s.dstPath)
	if err != nil {
		return err
	}
	defer destination.Close() // nolint
	_, err = io.Copy(destination, source)
	return err
}

func (s *filePair) stripModeline(lnum uint8) string {
	inF, err := os.Open(s.srcPath)
	errutils.Epanicf(err, "Can not open '%s': %s", s.srcPath, err) // nolint
	defer inF.Close()                                              // nolint

	tmpdir := os.Getenv("TMPDIR")
	outF, err := ioutil.TempFile(tmpdir, "prjstart-")
	errutils.Epanicf(err, "Can not create tempfile: %s", err) // nolint
	defer outF.Close()                                        // nolint

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
		errutils.Epanicf(err, "Error writing to file '%s': %s", outF.Name(), err) // nolint
	}
	return outF.Name()
}

func (s *filePair) render(mline uint8) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove modeline
	tempPath := s.stripModeline(mline)
	defer func() {
		os.Remove(tempPath)
	}()

	template.File2File(tempPath, s.dstPath, globals.Vars)
	return nil
}

// hasModeLine scans the first mlen lines for a modeline
// returns MLnone, MLrender depending on defined action
func (s *filePair) hasModeLine() (action int, lnum uint8) {
	action = MLnone
	len := s.mlen
	if len == 0 {
		len = 5
	}
	mlactions := hasML{}.Init()
	source, err := os.Open(s.srcPath)
	errutils.Efatalf(err, "Can not open file %s: %v", s.srcPath, err)

	defer source.Close()
	scner := bufio.NewScanner(source)
	for scner.Scan() {
		lnum++
		line := scner.Bytes()
		for _, mlaction := range mlactions {
			hasMatch := mlaction.Regex.Match(line)
			if hasMatch {
				action = mlaction.Action
				break
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
	ml = append(ml, regexCompile(`prj:ignore\W?`, MLignore))
	return ml
}

func regexCompile(rex string, action int) hasMLAction {
	regex, err := regexp.Compile(rex)
	errutils.Epanicf(err, "Error compiling regex: %s", err) // nolint
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
func renderDir(path string, prjvars *template.TmplVars) string {
	regex := regexp.MustCompile(`{{[^}}]+}}`)
	if !regex.MatchString(path) {
		return path
	}
	path = template.Txt2String(path, prjvars)
	return path
}
