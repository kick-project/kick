package gitclient

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	plumb "github.com/kick-project/kick/internal/resources/gitclient/plumbing"
	"github.com/kick-project/kick/internal/utils/errutils"
)

// Get Downloads using the data provider by getter
func Get(url string, g *plumb.Plumbing) (path string, err error) {
	err = g.Handler(url)
	if err != nil {
		return "", err
	}
	if g.Method() == plumb.NOOP {
		return g.Path(), nil
	}
	if g.Method() == plumb.SYNC {
		c := Gitclient{
			URL:   g.URL(),
			Local: g.Path(),
			Ref:   g.Branch(),
		}
		c.Sync()
	}
	return g.Path(), err
}

// Gitclient gitclient
type Gitclient struct {
	URL    string
	Local  string
	Output *os.File
	Ref    string
}

// Sync will download/synchronize with the upstream git repo
func (d *Gitclient) Sync() {
	d.Clone()
	d.Pull()
	d.Checkout(d.Ref)
}

// SetRef sets the default reference. See Checkout
func (d *Gitclient) SetRef(ref string) {
	d.Ref = ref
}

// Clone will clone a remote repository
func (d *Gitclient) Clone() {
	p := d.Local
	fmt.Printf("action=clone src=%s dest=%s\n", d.URL, p)
	if p == "" {
		return
	}
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		_, err := git.PlainClone(d.Local, false, &git.CloneOptions{
			URL:      d.URL,
			Progress: os.Stdout,
		})
		errutils.Efatalf("Can not clone %s: %v", d.URL, err)
	}
}

// Pull will pull from the remote repository
func (d *Gitclient) Pull() {
	p := d.Local
	if p == "" {
		return
	}
	r, err := git.PlainOpen(p)
	errutils.Elogf("Error opening path '%s': %+v", p, err)

	w, err := r.Worktree()
	errutils.Elogf("Error reading path '%s': %+v", p, err)

	pullopts := &git.PullOptions{}
	err = w.Pull(pullopts)
	if err != git.NoErrAlreadyUpToDate {
		errutils.Elogf("Error cloning %s: %+v", d.URL, err)
	}
}

// Tags will list all tags
func (d *Gitclient) Tags() []string {
	taglist := []string{}

	r := d.plainopen()
	iter, err := r.Tags()
	errutils.Elogf("Error listing tags %s: %+v", d.URL, err)

	fn := func(tag *plumbing.Reference) error {
		taglist = append(taglist, tag.Name().Short())
		return nil
	}
	err = iter.ForEach(fn)
	errutils.Elogf("Error listing tags %s: %+v", d.URL, err)

	return taglist
}

// Checkout checks out a reference. If ref is an empty string will checkout using the internally set ref
func (d *Gitclient) Checkout(ref string) {
	if ref != "" {
		d.Ref = ref
	}

	if d.Ref == "" {
		return
	}
	r, err := git.PlainOpen(d.Local)
	errutils.Epanicf("Error opening path '%s': %+v", d.Local, err)

	refObj, err := r.Reference(plumbing.ReferenceName(d.Ref), true)
	errutils.Epanicf("Error reading reference '%s' for path '%s': %+v", d.Ref, d.Local, err)

	chkops := &git.CheckoutOptions{
		Hash: refObj.Hash(),
	}

	w, err := r.Worktree()
	errutils.Epanicf("Error reading path '%s': %+v", d.Local, err)

	err = w.Checkout(chkops)
	errutils.Epanicf("Error checkout out: %+v", err)
}

func (d *Gitclient) plainopen() *git.Repository {
	p := d.Local
	if p == "" {
		return nil
	}
	r, err := git.PlainOpen(p)
	errutils.Elogf("Error opening path '%s': %+v", p, err)
	return r
}
