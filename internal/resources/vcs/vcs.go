package gitclient

import (
	"os"

	plumb "github.com/crosseyed/prjstart/internal/gitclient/plumbing"
	"github.com/crosseyed/prjstart/internal/utils/errutils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Get Downloads using the data provider by getter
func Get(url string, p *plumb.Plumbing) (path string, err error) {
	err = p.Handler(url)
	if err != nil {
		return "", err
	}
	if p.Method() == plumb.NOOP {
		return p.Path(), nil
	}
	if p.Method() == plumb.SYNC {
		vcs := New(Options{
			URL:  p.URL(),
			Path: p.Path(),
			Ref:  p.Branch(),
		})
		vcs.Sync()
	}
	return p.Path(), err
}

// VCS Version control system
type (
	// VCS the version contorl system
	VCS struct {
		uRL    string
		local  string
		stdout *os.File
		stderr *os.File
		ref    string
	}

	// Options provide options to New
	Options struct {
		URL    string
		Path   string
		Stdout *os.File
		Stderr *os.File
		Ref    string
	}
)

// New create a new VCS
func New(opts Options) *VCS {
	stdout := os.Stdout
	if opts.Stdout != nil {
		stdout = opts.Stdout
	}
	stderr := os.Stderr
	if opts.Stderr != nil {
		stderr = opts.Stderr
	}
	v := &VCS{
		uRL:    opts.URL,
		local:  opts.Path,
		stdout: stdout,
		stderr: stderr,
		ref:    opts.Ref,
	}
	return v
}

// Sync will download/synchronize with the upstream git repo
func (v *VCS) Sync() {
	v.Clone()
	v.Pull()
	v.Checkout(v.ref)
}

// SetRef sets the default reference. See Checkout
func (v *VCS) SetRef(ref string) {
	v.ref = ref
}

// Clone will clone a remote repository
func (v *VCS) Clone() {
	p := v.local
	if p == "" {
		return
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		_, err := git.PlainClone(v.local, false, &git.CloneOptions{
			URL:      v.uRL,
			Progress: v.stdout,
		})
		errutils.Efatalf("Can not clone %s: %v", v.uRL, err)
	}
}

// Pull will pull from the remote repository
func (v *VCS) Pull() {
	p := v.local
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
		errutils.Elogf("Error cloning %s: %+v", v.uRL, err)
	}
}

// Tags will list all tags
func (v *VCS) Tags() []string {
	taglist := []string{}

	r := v.plainopen()
	iter, err := r.Tags()
	errutils.Elogf("Error listing tags %s: %+v", v.uRL, err)

	fn := func(tag *plumbing.Reference) error {
		taglist = append(taglist, tag.Name().Short())
		return nil
	}
	err = iter.ForEach(fn)
	errutils.Elogf("Error listing tags %s: %+v", v.uRL, err)

	return taglist
}

// Checkout checks out a reference. If ref is an empty string will checkout using the internally set ref
func (v *VCS) Checkout(ref string) {
	if ref != "" {
		v.ref = ref
	}

	if v.ref == "" {
		return
	}
	r, err := git.PlainOpen(v.local)
	errutils.Epanicf("Error opening path '%s': %+v", v.local, err)

	refObj, err := r.Reference(plumbing.ReferenceName(v.ref), true)
	errutils.Epanicf("Error reading reference '%s' for path '%s': %+v", v.ref, v.local, err)

	chkops := &git.CheckoutOptions{
		Hash: refObj.Hash(),
	}

	w, err := r.Worktree()
	errutils.Epanicf("Error reading path '%s': %+v", v.local, err)

	err = w.Checkout(chkops)
	errutils.Epanicf("Error checkout out: %+v", err)
}

func (v *VCS) plainopen() *git.Repository {
	p := v.local
	if p == "" {
		return nil
	}
	r, err := git.PlainOpen(p)
	errutils.Elogf("Error opening path '%s': %+v", p, err)
	return r
}
