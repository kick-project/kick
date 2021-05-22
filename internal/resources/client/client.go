package client

import (
	"io"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-playground/validator"
	"github.com/kick-project/kick/internal/resources/client/plumb"
	"github.com/kick-project/kick/internal/resources/errs"
)

// Client
type Client struct {
	err    errs.HandlerIface
	stdout io.Writer
}

// Options for New function
type Options struct {
	Err    errs.HandlerIface `validate:"required"`
	Stdout io.Writer         `validate:"required"`
}

// New Client constructor
func New(opts *Options) *Client {
	v := validator.New()
	err := v.Struct(opts)
	if err != nil {
		panic(err)
	}
	return &Client{
		err:    opts.Err,
		stdout: opts.Stdout,
	}
}

// Get get url and clone/sync to path using ref.
// Defaults to default branch if ref is nil.
func (c *Client) Get(url, path, ref string) error {
	return nil
}

// GetPlumb same as get but url, path and ref are fetch from plumb.Plumb
func (c *Client) GetPlumb(p *plumb.Plumb) error {
	return c.Get(p.URL(), p.Path(), p.Ref())
}

// Sync will download/synchronize with the upstream git repo
func (d *Client) Sync(url, path, ref string) {
	d.Clone(url, path)
	d.Pull(path)
	d.Checkout(path, ref)
}

// Clone will clone a remote repository
func (d *Client) Clone(url, path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		errs.FatalF("Can not clone %s: %v", url, err)
	}
}

// Pull will pull from the remote repository
func (d *Client) Pull(path string) {
	r, err := git.PlainOpen(path)
	errs.LogF("Error opening path '%s': %+v", path, err)

	w, err := r.Worktree()
	errs.LogF("Error reading path '%s': %+v", path, err)

	pullopts := &git.PullOptions{}
	err = w.Pull(pullopts)
	if err != git.NoErrAlreadyUpToDate {
		errs.LogF("Error pulling %s: %+v", path, err)
	}
}

// Tags will list all tags
func (d *Client) Tags(path string) []string {
	taglist := []string{}

	r := d.plainopen(path)
	iter, err := r.Tags()
	errs.LogF("Error listing tags %s: %+v", path, err)

	fn := func(tag *plumbing.Reference) error {
		taglist = append(taglist, tag.Name().Short())
		return nil
	}
	err = iter.ForEach(fn)
	errs.LogF("Error listing tags %s: %+v", path, err)

	return taglist
}

// Checkout checks out a reference. If ref is an empty string will checkout using the internally set ref
func (d *Client) Checkout(path, ref string) {
	if ref == "" {
		return
	}
	r, err := git.PlainOpen(path)
	errs.PanicF("Error opening path '%s': %+v", path, err)

	refObj, err := r.Reference(plumbing.ReferenceName(ref), true)
	errs.PanicF("Error reading reference '%s' for path '%s': %+v", ref, path, err)

	chkops := &git.CheckoutOptions{
		Hash: refObj.Hash(),
	}

	w, err := r.Worktree()
	errs.PanicF("Error reading path '%s': %+v", path, err)

	err = w.Checkout(chkops)
	errs.PanicF("Error checkout out: %+v", err)
}

func (d *Client) plainopen(path string) *git.Repository {
	if path == "" {
		return nil
	}
	r, err := git.PlainOpen(path)
	errs.LogF("Error opening path '%s': %+v", path, err)
	return r
}
