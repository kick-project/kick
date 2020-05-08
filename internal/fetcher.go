package internal

import (
	"github.com/crosseyed/prjstart/internal/gitclient"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Fetcher struct {
	config *ConfigStruct
}

type RemoteTmpls struct {
	uri  string
	desc string
}

type parseFunc = func(uri string) (server string, path string, project string, match bool)

func NewFetcher(config *ConfigStruct) *Fetcher {
	return &Fetcher{
		config: config,
	}
}

// GetTmpl fetches the template using git and returns the local path or returns path if a local one
func (d *Fetcher) GetTmpl(tmpl string) string {
	for _, t := range d.config.Templates {
		if t.Name == tmpl {
			path := expandPath(t.URL)
			if strings.HasPrefix(path, "/") {
				return path
			}
			client := gitclient.New(t.URL, BaseProjectPath(d.config.Home), os.Stdout)
			if !client.LocalOnly() {
				client.Sync()
			}
			return client.LocalPath()
		}
	}
	return ""
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, path[2:])
	}
	return path
}

// GetAllSets fetches all sets
func (d *Fetcher) GetAllSets() []string {
	localSets := []string{}
	for _, url := range d.config.SetURLs {
		p := d.GetSet(url)
		localSets = append(localSets, p)
	}
	return localSets
}

// GetSet fetches sets and returns the patch
func (d *Fetcher) GetSet(uri string) string {
	client := gitclient.New(uri, BaseSetPath(d.config.Home), os.Stdout)
	if !client.LocalOnly() {
		client.Sync()
	}
	return client.LocalPath()
}