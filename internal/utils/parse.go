package utils

import (
	"fmt"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

// URLx take a URI or Something Specific to this project and break it into parts
type URLx struct {
	URL     string
	Scheme  string
	Project string
	Host    string
	Path    string
}

// Parse parses a URL or the link sets its internal attributes.
func (ux *URLx) Parse(url string) error {
	for _, parseF := range []parseFunc{httpParse, gitParse, sshParse, fileParse} {
		var match bool
		scheme, host, path, project, match := parseF(url)
		if match {
			ux.URL = url
			ux.Scheme = scheme
			ux.Project = project
			ux.Host = host
			ux.Path = path
			return nil
		}
	}
	err := fmt.Errorf("Could not parse url %s", url)
	return err
}

// Parse URL string
type parseFunc = func(uri string) (scheme, host, path, project string, match bool)

// Parse parses a URL or the like and returns a URLx pointer or nil if
// the parser failed to match a string.
func Parse(uri string) *URLx {
	u := &URLx{}
	u.Parse(uri)
	return u
}

func httpParse(uri string) (scheme, host, path, project string, match bool) {
	r := regexp.MustCompile(`^(https?)://([^(?:/|:)]+)(?:/|:\d+)(.*?)([^/]+?)(?:\.git)?$`)
	m := r.FindStringSubmatch(uri)
	if len(m) > 3 {
		return m[1], m[2], filepath.Clean(filepath.Join(m[2], m[3])), m[4], true
	}
	return "", "", "", "", false
}

func gitParse(uri string) (scheme, host, path, project string, match bool) {
	r := regexp.MustCompile(`^(git)@([^(?:/|:)]+)(?:/|:)(.*?)([^/]+?)(?:\.git)?$`)
	m := r.FindStringSubmatch(uri)
	if len(m) > 3 {
		return m[1], m[2], filepath.Clean(filepath.Join(m[2], m[3])), m[4], true
	}
	return "", "", "", "", false
}

func sshParse(uri string) (scheme, host, path, project string, match bool) {
	r := regexp.MustCompile(`^(ssh)://([^(?:/|:)]+)(?:/|:\d+)(.*?)([^/]+?)(?:\.git)?$`)
	m := r.FindStringSubmatch(uri)
	if len(m) > 3 {
		return m[1], m[2], filepath.Clean(filepath.Join(m[2], m[3])), m[4], true
	}
	return "", "", "", "", false
}

func fileParse(uri string) (scheme, host, path, project string, match bool) {
	r := regexp.MustCompile(`^(file)://(/.*?)([^/]+?)/?$`)
	m := r.FindStringSubmatch(uri)
	if len(m) > 1 {
		return m[1], "", filepath.Clean(filepath.Join(m[2], m[3])), m[3], true
	}
	return "", "", "", "", false
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, path[2:])
	}
	return path
}
