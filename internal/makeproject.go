package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/crosseyed/prjstart/internal/template"
	"github.com/crosseyed/prjstart/internal/utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	MLnone = iota
	MLrender
)

type MakeProject struct {
	localpath string
	src       string
	dest      string
	tmp       string
}

func (s *MakeProject) SetTemp(id string) {
	d, err := ioutil.TempDir(os.Getenv("TEMP"), fmt.Sprintf("prjstart-%s-", id))
	utils.ChkErr(err, utils.Epanicf)
	s.tmp = d
}

// SetSrc sets the source template and localpath
func (s *MakeProject) SetSrc(src string) {
	dwnlder := NewFetcher(Config)
	localpath := dwnlder.GetTmpl(src)
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

func (s *MakeProject) SetDest(dest string) {
	s.dest = dest
}

func (s *MakeProject) Run() int {
	path := s.localpath
	base := s.localpath
	skipRegex, err := regexp.Compile(fmt.Sprintf(`^%s/.git(?:/|$)`, base))
	utils.ChkErr(err, utils.Epanicf)
	s.checkDstExists()
	errWalk := filepath.Walk(path, func(srcPath string, info os.FileInfo, err error) error {
		if skipRegex.MatchString(srcPath) {
			return nil
		}
		relative := strings.Replace(srcPath, base, "", 1)
		relative = RenderDir(relative, Vars)
		dstPath := filepath.Join(s.tmp, relative)

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
		err = pair.Route()
		utils.ChkErr(err, utils.Epanicf)

		return nil
	})
	if errWalk != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Abort creating project: %s", errWalk.Error())
		return 255
	}

	err = os.Rename(s.tmp, s.dest)
	utils.ChkErr(err, utils.Epanicf)
	return 0
}

func (s *MakeProject) checkDstExists() {
	stat, err := os.Stat(s.dest)
	if !os.IsNotExist(err) {
		utils.ChkErr(err, utils.Epanicf)
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

func (s *filePair) Route() error {
	action, lnum := s.hasModeLine()
	switch {
	case s.srcInfo.IsDir():
		return s.Mkdir()
	case lnum > 0 && action == MLrender:
		s.Render(lnum)
	case s.srcInfo.Mode().IsRegular():
		return s.Copy()
	default:
		msg := fmt.Sprintf("error FILENOTREGULAR: %s\n", s.dstPath)
		fmt.Println(msg)
		return errors.New(msg)
	}
	return nil
}

func (s *filePair) Mkdir() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := os.Stat(s.dstPath); os.IsNotExist(err) {
		err = os.Mkdir(s.dstPath, 0777)
		utils.ChkErr(err, utils.Epanicf)
	}
	return nil
}

func (s *filePair) Copy() error {
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

func (s *filePair) StripModeline(lnum uint8) string {
	inF, err := os.Open(s.srcPath)
	utils.ChkErr(err, utils.Epanicf, "Can not open '%s': %s", s.srcPath, err) // nolint
	defer inF.Close()                                                         // nolint

	tmpdir := os.Getenv("TMPDIR")
	outF, err := ioutil.TempFile(tmpdir, "prjstart-")
	utils.ChkErr(err, utils.Epanicf, "Can not create tempfile: %s", err) // nolint
	defer outF.Close()                                                   // nolint

	var cnt uint8
	cnt = 0
	scner := bufio.NewScanner(inF)
	scner.Split(ScanLines)
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
		utils.ChkErr(err, utils.Epanicf, "Error writing to file '%s': %s", outF.Name(), err) // nolint
	}
	return outF.Name()
}

func (s *filePair) Render(mline uint8) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove modeline
	tempPath := s.StripModeline(mline)
	defer func() {
		os.Remove(tempPath)
	}()

	template.TmplFile2File(tempPath, s.dstPath, Vars)
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
	regex, err := regexp.Compile(`prj:(render)`)
	utils.ChkErr(err, utils.Epanicf, "Error compiling regex: %s", err) // nolint
	source, err := os.Open(s.srcPath)
	utils.ChkErr(err, utils.Epanicf, "Error opening read file '%s': %s", s.srcPath, err)
	defer source.Close()

	scner := bufio.NewScanner(source)
	for scner.Scan() {
		lnum++
		line := scner.Bytes()
		hasMatch := regex.Match(line)
		if hasMatch {
			action = MLrender
			break
		}
		if lnum > len {
			lnum = 0
			break
		}
	}
	return action, lnum
}

// ScanLines is a split function for a Scanner that returns each line of
// text, WITHOUT stripping any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r?\n`.
// The last non-empty line of input will be returned even if it has no
// newline.
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
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

// RenderDir scans directory names for template markers and renders the directory path as a template
func RenderDir(path string, prjvars *template.TmplVars) string {
	regex := regexp.MustCompile(`{{[^}}]+}}`)
	if !regex.MatchString(path) {
		return path
	}
	path = template.TmplTxt2String(path, prjvars)
	return path
}

// BaseProjectPath returns the path to the project path
func BaseProjectPath(home string) string {
	if home == "" {
		home = os.Getenv("HOME")
	}
	p := filepath.Join(home, ".prjstart", "projects")
	return p
}

func BaseSetPath(home string) string {
	if home == "" {
		home = os.Getenv("HOME")
	}
	p := filepath.Join(home, ".prjstart", "sets")
	return p
}
