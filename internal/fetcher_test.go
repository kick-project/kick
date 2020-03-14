package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//
// Tests
//
func TestGetTmpl(t *testing.T) {
	dwnlder := NewFetcher(Config)
	lp := dwnlder.GetTmpl("tmpl")
	assert.DirExists(t, lp)
}

func TestFetcher_GetAllSets(t *testing.T) {
	//localSets := TestServer.Fetcher.GetAllSets()
	//for _, localSet := range localSets {
	//	assert.DirExists(t, localSet)
	//}
}

func TestHttpParseHTTP(t *testing.T) {
	url := "http://git.com/serve/git/template.git"
	expectedServer := "git.com"
	expectedProject := "template"
	expectedPath := "git.com/serve/git"

	testParsing(t, httpParse, url, expectedServer, expectedPath, expectedProject)
}

func TestHttpParseHTTPS(t *testing.T) {
	url := "https://github.com/serve/git/gotmpl.git"
	expectedServer := "github.com"
	expectedProject := "gotmpl"
	expectedPath := "github.com/serve/git"

	testParsing(t, httpParse, url, expectedServer, expectedPath, expectedProject)
}

func TestHttpParseHTTPPORT(t *testing.T) {
	url := "http://bitbucket.org:5000/serve/git/pytmpl.git"
	expectedServer := "bitbucket.org"
	expectedProject := "pytmpl"
	expectedPath := "bitbucket.org/serve/git"

	testParsing(t, httpParse, url, expectedServer, expectedPath, expectedProject)
}

func TestHttpParseHTTPSPORT(t *testing.T) {
	url := "https://127.0.0.1:5000/serve/git/template.git"
	expectedServer := "127.0.0.1"
	expectedProject := "template"
	expectedPath := "127.0.0.1/serve/git"

	testParsing(t, httpParse, url, expectedServer, expectedPath, expectedProject)
}

func TestGitParse(t *testing.T) {
	url := "git@127.0.0.1/serve/git/template.git"
	expectedServer := "127.0.0.1"
	expectedProject := "template"
	expectedPath := "127.0.0.1/serve/git"

	testParsing(t, gitParse, url, expectedServer, expectedPath, expectedProject)
}

func TestGitParse2(t *testing.T) {
	url := "git@127.0.0.1:serve/git/template.git"
	expectedServer := "127.0.0.1"
	expectedProject := "template"
	expectedPath := "127.0.0.1/serve/git"

	testParsing(t, gitParse, url, expectedServer, expectedPath, expectedProject)
}

func TestSshParse(t *testing.T) {
	url := "ssh://github.com/user/tmpl.git"
	expectedServer := "github.com"
	expectedProject := "tmpl"
	expectedPath := "github.com/user"

	testParsing(t, sshParse, url, expectedServer, expectedPath, expectedProject)
}

func TestFileParse(t *testing.T) {
	url := "file:///home/user/workspace/tmpl"
	expectedServer := "::local::"
	expectedProject := "tmpl"
	expectedPath := "/home/user/workspace/tmpl" // Locals have a special path

	testParsing(t, fileParse, url, expectedServer, expectedPath, expectedProject)
}

func testParsing(t *testing.T, f parseFunc, url string, expectedServer string, expectedPath string, expectedProject string) {
	server, path, project, match := f(url)
	assert.True(t, match)
	assert.Equal(t, expectedServer, server)
	assert.Equal(t, expectedProject, project)
	assert.Equal(t, expectedPath, path)
	server, path, project = parseGitRemote(url)
	assert.Equal(t, expectedServer, server)
	assert.Equal(t, expectedProject, project)
	assert.Equal(t, expectedPath, path)
}
