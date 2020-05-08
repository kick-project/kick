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
	uri     string
	basedir string
	output  *os.File
	ref     string
}

func New(uri, base string, output *os.File) *Gitclient {
	s := &Gitclient{
		uri:     uri,
		basedir: base,
		output:  output,
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
