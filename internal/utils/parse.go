package utils

import (
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

// Parse URL string
type parseFunc = func(uri string) (server string, path string, project string, match bool)

func ParseGitRemote(uri string) (server, path, project string) {
	for _, parseF := range []parseFunc{httpParse, gitParse, sshParse, fileParse} {
		var match bool
		server, path, project, match = parseF(uri)
		if match {
			return server, path, project
		}
	}
	return "", "", ""
}

func httpParse(uri string) (server string, path string, project string, match bool) {
	r := regexp.MustCompile(`^https?://([^(?:/|:)]+)(?:/|:\d+)(.*?)([^/]+?)(?:\.git)?$`)
	m := r.FindStringSubmatch(uri)
	if len(m) > 3 {
		return m[1], filepath.Clean(filepath.Join(m[1], m[2])), m[3], true
	}
	return "", "", "", false
}

func gitParse(uri string) (server, path, project string, match bool) {
	r := regexp.MustCompile(`^git@([^(?:/|:)]+)(?:/|:)(.*?)([^/]+?)(?:\.git)?$`)
	m := r.FindStringSubmatch(uri)
	if len(m) > 3 {
		return m[1], filepath.Clean(filepath.Join(m[1], m[2])), m[3], true
	}
	return "", "", "", false
}

func sshParse(uri string) (server, path, project string, match bool) {
	r := regexp.MustCompile(`^ssh://([^(?:/|:)]+)(?:/|:\d+)(.*?)([^/]+?)(?:\.git)?$`)
	m := r.FindStringSubmatch(uri)
	if len(m) > 3 {
		return m[1], filepath.Clean(filepath.Join(m[1], m[2])), m[3], true
	}
	return "", "", "", false
}

func fileParse(uri string) (server, path, project string, match bool) {
	r := regexp.MustCompile(`^(?:file://)?(/.*?)([^/]+?)/?$`)
	m := r.FindStringSubmatch(uri)
	if len(m) > 1 {
		return "::local::", filepath.Clean(filepath.Join(m[1], m[2])), m[2], true
	}
	return "", "", "", false
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, path[2:])
	}
	return path
}
