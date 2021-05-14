package gitclient

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/kick-project/kick/internal/resources/errs"
	plumb "github.com/kick-project/kick/internal/resources/gitclient/plumbing"
)

// Get Downloads using the data provider by getter
func Get(url string, g plumb.PlumbingIface) (path string, err error) {
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
	URL    string   `copier:"must"`
	Local  string   `copier:"must"`
	Output *os.File `copier:"must"`
	Ref    string   `copier:"must"`
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
	if p == "" {
		return
	}
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		_, err := git.PlainClone(d.Local, false, &git.CloneOptions{
			URL:      d.URL,
			Progress: os.Stdout,
		})
		errs.FatalF("Can not clone %s: %v", d.URL, err)
	}
}

// Pull will pull from the remote repository
func (d *Gitclient) Pull() {
	p := d.Local
	if p == "" {
		return
	}
	r, err := git.PlainOpen(p)
	errs.LogF("Error opening path '%s': %+v", p, err)

	w, err := r.Worktree()
	errs.LogF("Error reading path '%s': %+v", p, err)

	pullopts := &git.PullOptions{}
	err = w.Pull(pullopts)
	if err != git.NoErrAlreadyUpToDate {
		errs.LogF("Error cloning %s: %+v", d.URL, err)
	}
}

// Tags will list all tags
func (d *Gitclient) Tags() []string {
	taglist := []string{}

	r := d.plainopen()
	iter, err := r.Tags()
	errs.LogF("Error listing tags %s: %+v", d.URL, err)

	fn := func(tag *plumbing.Reference) error {
		taglist = append(taglist, tag.Name().Short())
		return nil
	}
	err = iter.ForEach(fn)
	errs.LogF("Error listing tags %s: %+v", d.URL, err)

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
	errs.PanicF("Error opening path '%s': %+v", d.Local, err)

	refObj, err := r.Reference(plumbing.ReferenceName(d.Ref), true)
	errs.PanicF("Error reading reference '%s' for path '%s': %+v", d.Ref, d.Local, err)

	chkops := &git.CheckoutOptions{
		Hash: refObj.Hash(),
	}

	w, err := r.Worktree()
	errs.PanicF("Error reading path '%s': %+v", d.Local, err)

	err = w.Checkout(chkops)
	errs.PanicF("Error checkout out: %+v", err)
}

func (d *Gitclient) plainopen() *git.Repository {
	p := d.Local
	if p == "" {
		return nil
	}
	r, err := git.PlainOpen(p)
	errs.LogF("Error opening path '%s': %+v", p, err)
	return r
}
