package mclient

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/crosseyed/prjstart/internal/config"
	"github.com/crosseyed/prjstart/internal/gitclient"
	"github.com/crosseyed/prjstart/internal/utils"
)

// Metadata is the metadata client. It pulls from git repositories
type Metadata struct {
	Config *config.Config
}

// Get fetches the metadata using git and returns the local path or returns the path if local.
func (m *Metadata) Get(url string) (lpath string) {
	url = utils.ExpandPath(url)
	if strings.HasPrefix(url, "/") {
		return url
	}
	if !localOnly(url) {
		lpath := m.localPath(url)
		client := gitclient.Gitclient{
			URL:    url,
			Local:  lpath,
			Output: os.Stdout,
		}
		client.Sync()
		return lpath
	}
	return ""
}

func (m *Metadata) EachTag(fn func(string) (stop bool), url string) {
	url = utils.ExpandPath(url)
	if !localOnly(url) {
		lpath := m.localPath(url)
		client := gitclient.Gitclient{
			URL:    url,
			Local:  lpath,
			Output: os.Stdout,
		}
		client.Sync()
		client.EachTag(fn)
	}
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
func (m *Metadata) localPath(url string) string {
	server, srvPath, dir := utils.ParseGitRemote(url)
	if server == "::local::" {
		return srvPath
	}
	if srvPath == "" || dir == "" {
		return ""
	}
	base := utils.BaseMetadataPath(m.Config.Home)
	p := filepath.Join(base, srvPath, dir)
	return p
}
