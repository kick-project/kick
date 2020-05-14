package tclient

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/crosseyed/prjstart/internal/config"
	"github.com/crosseyed/prjstart/internal/gitclient"
	"github.com/crosseyed/prjstart/internal/utils"
)

// Parse URL string
type parseFunc = func(uri string) (server string, path string, project string, match bool)

// TClient is the template client
type TClient struct {
	Config *config.Config
}

// Get fetches the template using git and returns the local path or returns path if a local one
func (c *TClient) Get(tmpl string) string {
	for _, t := range c.Config.Templates {
		if t.Name == tmpl {
			url := utils.ExpandPath(t.URL)
			if strings.HasPrefix(url, "/") {
				return url
			}
			if !localOnly(t.URL) {
				lpath := c.localPath(url)
				client := gitclient.Gitclient{
					URL:    t.URL,
					Local:  lpath,
					Output: os.Stdout,
				}
				client.Sync()
				return lpath
			}
		}
	}
	return ""
}

// localOnly determines that the url is a local path
func localOnly(url string) bool {
	server, _, _ := utils.ParseGitRemote(url)
	if server == "::local::" {
		return true
	}
	return false
}

// localPath determines local path to template
func (c *TClient) localPath(url string) string {
	server, srvPath, dir := utils.ParseGitRemote(url)
	if server == "::local::" {
		return srvPath
	}
	if srvPath == "" || dir == "" {
		return ""
	}
	base := utils.BaseProjectPath(c.Config.Home)
	p := filepath.Join(base, srvPath, dir)
	return p
}
