package internal

import (
	"github.com/crosseyed/prjstart/internal/gitclient"
	"github.com/crosseyed/prjstart/internal/utils"
	"os"
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
			path := utils.ExpandPath(t.URL)
			if strings.HasPrefix(path, "/") {
				return path
			}
			options := gitclient.Options {
				Uri:     t.URL,
				BaseDir: BaseProjectPath(d.config.Home),
				OutPut:  os.Stdout,
			}
			client := gitclient.New(options)
			if !client.LocalOnly() {
				client.Sync()
			}
			return client.LocalPath()
		}
	}
	return ""
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
	options := gitclient.Options {
		Uri:     uri,
		BaseDir: BaseProjectPath(d.config.Home),
		OutPut:  os.Stdout,
	}
	client := gitclient.New(options)
	if !client.LocalOnly() {
		client.Sync()
	}
	return client.LocalPath()
}