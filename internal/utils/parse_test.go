package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// Tests
//
func TestHttpParseHTTP(t *testing.T) {
	url := "http://git.com/serve/git/template.git"
	expectedScheme := "http"
	expectedServer := "git.com"
	expectedProject := "template"
	expectedPath := "git.com/serve/git"

	testParsing(t, httpParse, url, expectedScheme, expectedServer, expectedPath, expectedProject)
}

func TestHttpParseHTTPS(t *testing.T) {
	url := "https://github.com/serve/git/gotmpl.git"
	expectedScheme := "https"
	expectedServer := "github.com"
	expectedProject := "gotmpl"
	expectedPath := "github.com/serve/git"

	testParsing(t, httpParse, url, expectedScheme, expectedServer, expectedPath, expectedProject)
}

func TestHttpParseHTTPPORT(t *testing.T) {
	url := "http://bitbucket.org:5000/serve/git/pytmpl.git"
	expectedScheme := "http"
	expectedServer := "bitbucket.org"
	expectedProject := "pytmpl"
	expectedPath := "bitbucket.org/serve/git"

	testParsing(t, httpParse, url, expectedScheme, expectedServer, expectedPath, expectedProject)
}

func TestHttpParseHTTPSPORT(t *testing.T) {
	url := "https://127.0.0.1:5000/serve/git/template.git"
	expectedScheme := "https"
	expectedServer := "127.0.0.1"
	expectedProject := "template"
	expectedPath := "127.0.0.1/serve/git"

	testParsing(t, httpParse, url, expectedScheme, expectedServer, expectedPath, expectedProject)
}

func TestGitParse(t *testing.T) {
	url := "git@127.0.0.1/serve/git/template.git"
	expectedScheme := "git"
	expectedServer := "127.0.0.1"
	expectedProject := "template"
	expectedPath := "127.0.0.1/serve/git"

	testParsing(t, gitParse, url, expectedScheme, expectedServer, expectedPath, expectedProject)
}

func TestGitParse2(t *testing.T) {
	url := "git@127.0.0.1:serve/git/template.git"
	expectedScheme := "git"
	expectedServer := "127.0.0.1"
	expectedProject := "template"
	expectedPath := "127.0.0.1/serve/git"

	testParsing(t, gitParse, url, expectedScheme, expectedServer, expectedPath, expectedProject)
}

func TestSshParse(t *testing.T) {
	url := "ssh://github.com/user/tmpl.git"
	expectedScheme := "ssh"
	expectedServer := "github.com"
	expectedProject := "tmpl"
	expectedPath := "github.com/user"

	testParsing(t, sshParse, url, expectedScheme, expectedServer, expectedPath, expectedProject)
}

func TestFileParse(t *testing.T) {
	url := "file:///home/user/workspace/tmpl"
	expectedScheme := "file"
	expectedServer := ""
	expectedProject := "tmpl"
	expectedPath := "/home/user/workspace/tmpl" // Locals have a special path

	testParsing(t, fileParse, url, expectedScheme, expectedServer, expectedPath, expectedProject)
}

func testParsing(t *testing.T, f parseFunc, url, expectedScheme, expectedServer, expectedPath, expectedProject string) {
	scheme, host, path, project, match := f(url)
	assert.True(t, match)
	assert.Equal(t, expectedScheme, scheme)
	assert.Equal(t, expectedServer, host)
	assert.Equal(t, expectedProject, project)
	assert.Equal(t, expectedPath, path)
	u, _ := Parse(url)
	assert.Equal(t, expectedScheme, u.Scheme)
	assert.Equal(t, expectedServer, u.Host)
	assert.Equal(t, expectedProject, u.Project)
	assert.Equal(t, expectedPath, u.Path)
}
