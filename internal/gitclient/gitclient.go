package gitclient

import (
	"fmt"
	"os"

	"github.com/crosseyed/prjstart/internal/utils"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

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
	d.Checkout("")
}

func (d *Gitclient) EachTag(fn func(tag string) (stop bool)) {
	wd, err := os.Getwd()
	errutils.Epanicf(err, "Can not get working directory: %v", err)
	err = os.Chdir(wd)
	errutils.Epanicf(err, "Can not change directory to %s: %v", wd, err)

	defer os.Chdir(wd)

	err = os.Chdir(d.Local)
	errutils.Epanicf(err, "Can not change directory to %s: %v", wd, err)
	for _, t := range d.Tags() {
		d.Checkout(t)
		stop := fn(t)
		if stop {
			break
		}
	}
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
	if _, err := os.Stat(p); os.IsNotExist(err) {
		_, err := git.PlainClone(d.Local, false, &git.CloneOptions{
			URL:      d.URL,
			Progress: os.Stdout,
		})
		if err != nil {
			fmt.Printf("Can not clone %s: %s\n", d.URL, err.Error())
			utils.Exit(-1)
		}
	}
}

// Pull will pull from the remote repository
func (d *Gitclient) Pull() {
	p := d.Local
	if p == "" {
		return
	}
	r, err := git.PlainOpen(p)
	errutils.Elogf(err, "Error opening path '%s': %+v", p, err)

	w, err := r.Worktree()
	errutils.Elogf(err, "Error reading path '%s': %+v", p, err)

	pullopts := &git.PullOptions{}
	err = w.Pull(pullopts)
	if err != git.NoErrAlreadyUpToDate {
		errutils.Elogf(err, "Error cloning %s: %+v", d.URL, err)
	}
}

// Tags will list all tags
func (d *Gitclient) Tags() []string {
	taglist := []string{}

	r := d.plainopen()
	iter, err := r.Tags()
	errutils.Elogf(err, "Error listing tags %s: %+v", d.URL, err)

	fn := func(tag *plumbing.Reference) error {
		taglist = append(taglist, tag.Name().Short())
		return nil
	}
	err = iter.ForEach(fn)
	errutils.Elogf(err, "Error listing tags %s: %+v", d.URL, err)

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
	errutils.Elogf(err, "Error opening path '%s': %+v", d.Local, err)

	refObj, err := r.Reference(plumbing.ReferenceName(d.Ref), true)
	errutils.Elogf(err, "Error reading reference for path '%s': %+v", d.Local, err)

	chkops := &git.CheckoutOptions{
		Hash: refObj.Hash(),
	}

	w, err := r.Worktree()
	errutils.Elogf(err, "Error reading path '%s': %+v", d.Local, err)

	err = w.Checkout(chkops)
	errutils.Elogf(err, "Error checkout out: %+v", err)
}

func (d *Gitclient) plainopen() *git.Repository {
	p := d.Local
	if p == "" {
		return nil
	}
	r, err := git.PlainOpen(p)
	errutils.Elogf(err, "Error opening path '%s': %+v", p, err)
	return r
}
