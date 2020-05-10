package gitclient

import (
	"fmt"
	"github.com/crosseyed/prjstart/internal/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path/filepath"
	"regexp"
)

// Parse URL string
type parseFunc = func(uri string) (server string, path string, project string, match bool)

// Gitclient gitclient
type Gitclient struct {
	uri       string
	basedir   string
	forcepath bool
	output    *os.File
	ref       string
}

// Options is a set of options to past to New
type Options struct {
	// Uri is git git URI to check out
	Uri       string
	// BaseDir is the directory where the structure is built
	BaseDir   string 
	// ForcePath sets base directory 
	ForcePath bool
	// OutPut redirect output to file
	OutPut    *os.File
}

func New(options Options) *Gitclient {
	s := &Gitclient{
		uri:       options.Uri,
		basedir:   options.BaseDir,
		output:    options.OutPut,
		forcepath: options.ForcePath,
	}
	return s
}

// LocalOnly determines if a repo is local only
func (d *Gitclient) LocalOnly() bool {
	server, _, _ := parseGitRemote(d.uri)
	if server == "::local::" {
		return true
	}
	return false
}

// Sync will download/synchronize with the upstream git repo
func (d *Gitclient) Sync() {
	d.Clone()
	d.Pull()
	d.Checkout("")
}

// SetRef sets the default reference. See Checkout
func (d *Gitclient) SetRef(ref string) {
	d.ref = ref
}

// Clone will clone a remote repository
func (d *Gitclient) Clone() {
	p := d.LocalPath()
	if p == "" {
		return
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		_, err := git.PlainClone(d.LocalPath(), false, &git.CloneOptions{
			URL:      d.uri,
			Progress: os.Stdout,
		})
		if err != nil {
			fmt.Printf("Can not clone %s: %s\n", d.uri, err.Error())
			utils.Exit(-1)
		}
	}
}

// Pull will pull from the remote repository
func (d *Gitclient) Pull() {
	p := d.LocalPath()
	if p == "" {
		return
	}
	r, err := git.PlainOpen(p)
	utils.ChkErr(err, utils.Elogf, "Error opening path '%s': %+v", p, err)

	w, err := r.Worktree()
	utils.ChkErr(err, utils.Elogf, "Error reading path '%s': %+v", p, err)

	pullopts := &git.PullOptions{}
	err = w.Pull(pullopts)
	if err != git.NoErrAlreadyUpToDate {
		utils.ChkErr(err, utils.Elogf, "Error cloning %s: %+v", d.uri, err)
	}
}

// Tags will list all tags
func (d *Gitclient) Tags() []string {
	taglist := []string{}

	r := d.plainopen()
	iter, err := r.Tags()
	utils.ChkErr(err, utils.Elogf, "Error listing tags %s: %+v", d.uri, err)

	fn := func(tag *plumbing.Reference) error {
		taglist = append(taglist, tag.Name().Short())
		return nil
	}
	err = iter.ForEach(fn)
	utils.ChkErr(err, utils.Elogf, "Error listing tags %s: %+v", d.uri, err)

	return taglist
}

// Checkout checks out a reference. If ref is an empty string will checkout using the internally set ref
func (d *Gitclient) Checkout(ref string) {
	if ref != "" {
		d.ref = ref
	}

	if d.ref == "" {
		return
	}

	p := d.LocalPath()
	if p == "" {
		return
	}

	r, err := git.PlainOpen(p)
	utils.ChkErr(err, utils.Elogf, "Error opening path '%s': %+v", p, err)

	refObj, err := r.Reference(plumbing.ReferenceName(d.ref), true)
	utils.ChkErr(err, utils.Elogf, "Error reading reference for path '%s': %+v")

	chkops := &git.CheckoutOptions{
		Hash: refObj.Hash(),
	}

	w, err := r.Worktree()
	utils.ChkErr(err, utils.Elogf, "Error reading path '%s': %+v", p, err)

	err = w.Checkout(chkops)
	utils.ChkErr(err, utils.Elogf, "Error checkout out: %+v", err)
}

// LocalPath The local path
func (d *Gitclient) LocalPath() string {
	server, srvPath, dir := parseGitRemote(d.uri)
	if server == "::local::" {
		return srvPath
	}
	if srvPath == "" || dir == "" {
		return ""
	}
	p := filepath.Join(d.basedir, srvPath, dir)
	return p
}

func (d *Gitclient) plainopen() *git.Repository {
	p := d.LocalPath()
	if p == "" {
		return nil
	}
	r, err := git.PlainOpen(p)
	utils.ChkErr(err, utils.Elogf, "Error opening path '%s': %+v", p, err)
	return r
}

func parseGitRemote(uri string) (server, path, project string) {
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
