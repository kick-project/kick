package vcs

import (
	"fmt"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/kick-project/kick/internal/resources/errs"
)

// VCS package information
type VCS struct {
	err errs.HandlerIface
}

// Options constructor options
type Options struct {
	Err errs.HandlerIface
}

// New constructor
func New(opts *Options) *VCS {
	return &VCS{
		err: opts.Err,
	}
}

// Clone will clone a remote repository
func (i *VCS) Clone(url, path string) (repo *Repo, err error) {
	latest := false
	_, statErr := os.Stat(path)
	if os.IsNotExist(statErr) {
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
		})
		if i.err.LogF("Can not clone %s: %v", url, err) {
			return nil, fmt.Errorf("clone error: %w", err)
		}
		latest = true
	}

	repo, err = i.Open(path)
	if err != nil {
		return nil, fmt.Errorf(`error cloning %s, can not open %s: %w`, url, path, err)
	}
	if !latest {
		err = repo.Pull()
		if err != nil {
			return nil, fmt.Errorf(`error cloning %s, can not pull: %w`, url, err)
		}
	}
	return
}

func (i *VCS) Open(path string) (repo *Repo, err error) {
	repo = &Repo{
		path: path,
		err:  i.err,
	}
	r, err := git.PlainOpen(path)
	if i.err.LogF("Error opening path '%s': %+v", path, err) {
		return nil, fmt.Errorf("pull error: %w", err)
	}
	repo.repo = r
	return
}

type Repo struct {
	path string
	repo *git.Repository
	err  errs.HandlerIface
}

// Checkout checks out a reference
func (r *Repo) Checkout(ref string) error {
	if ref == "" {
		return nil
	}
	iter, err := r.repo.References()
	if err != nil {
		return fmt.Errorf(`checkout error: %w`, err)
	}
	var targetRef *plumbing.Reference
	fn := func(curRef *plumbing.Reference) error {
		if curRef.Name().Short() == ref {
			targetRef = curRef
		}
		return nil
	}
	err = iter.ForEach(fn)
	if err != nil {
		return fmt.Errorf(`reference error: %w`, err)
	}

	if targetRef == nil {
		return fmt.Errorf(`could not find reference "%s"`, ref)
	}

	refObj, err := r.repo.Reference(targetRef.Name(), true)
	if r.err.LogF("Error reading reference '%s' for path '%s': %+v", ref, r.path, err) {
		return fmt.Errorf("checkout error: %w", err)
	}

	chkops := &git.CheckoutOptions{
		Hash: refObj.Hash(),
	}

	w, err := r.repo.Worktree()
	if r.err.LogF("Error reading path '%s': %+v", r.path, err) {
		return fmt.Errorf("checkout error: %w", err)
	}

	err = w.Checkout(chkops)
	if r.err.LogF("Error checkout out: %+v", err) {
		return fmt.Errorf("checkout error: %w", err)
	}
	return nil
}

func (r *Repo) Pull() error {
	w, err := r.repo.Worktree()
	if r.err.LogF("Error reading path '%s': %+v", r.path, err) {
		return fmt.Errorf("pull error: %w", err)
	}

	remotes, err := r.repo.Remotes()
	if err != nil {
		return fmt.Errorf("listing remotes error: %w", err)
	}
	if len(remotes) == 0 {
		return nil
	}

	pullopts := &git.PullOptions{}
	err = w.Pull(pullopts)
	if err != git.NoErrAlreadyUpToDate {
		if r.err.LogF("Error pulling %s: %+v", r.path, err) {
			return fmt.Errorf("pull error: %w", err)
		}
	}

	return nil
}

// Tags will list all tags
func (r *Repo) Tags() (tags []string) {
	iter, err := r.repo.Tags()
	r.err.LogF("Error listing tags %s: %+v", r.path, err)

	fn := func(tag *plumbing.Reference) error {
		tags = append(tags, tag.Name().Short())
		return nil
	}
	err = iter.ForEach(fn)
	r.err.LogF("Error listing tags %s: %+v", r.path, err)

	return tags
}

// Versions will return a list of versions that match semantic versioning
func (r *Repo) Versions() (versions []string) {
	regex := regexp.MustCompile(`^\d+\.\d+(?:\.\d+)?$`)
	for _, t := range r.Tags() {
		if regex.MatchString(t) {
			versions = append(versions, t)
		}
	}
	return
}
