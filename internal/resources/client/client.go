package client

import (
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-playground/validator"
	"github.com/kick-project/kick/internal/di/callbacks"
	"github.com/kick-project/kick/internal/resources/client/plumb"
	"github.com/kick-project/kick/internal/resources/errs"
)

// Client
type Client struct {
	err            errs.HandlerIface
	stdout         io.Writer
	plumbRepos     callbacks.MakePlumb
	plumbTemplates callbacks.MakePlumb
}

// Options for New function
type Options struct {
	Err                errs.HandlerIface   `validate:"required"`
	Stdout             io.Writer           `validate:"required"`
	CallPlumbRepos     callbacks.MakePlumb `validate:"required"`
	CallPlumbTemplates callbacks.MakePlumb `validate:"required"`
}

// New Client constructor
func New(opts *Options) *Client {
	v := validator.New()
	err := v.Struct(opts)
	if err != nil {
		panic(err)
	}
	return &Client{
		err:            opts.Err,
		stdout:         opts.Stdout,
		plumbRepos:     opts.CallPlumbRepos,
		plumbTemplates: opts.CallPlumbTemplates,
	}
}

// Get get url and clone/sync to path using ref.
// Defaults to default branch if ref is nil.
func (c *Client) Get(url, path, ref string) error {
	err := c.Sync(url, path, ref)
	if err != nil {
		return fmt.Errorf("get error: %w", err)
	}
	return nil
}

// GetPlumb same as get but url, path and ref are fetch from plumb.Plumb
func (c *Client) GetPlumb(p *plumb.Plumb) error {
	switch p.Method() {
	case plumb.NOOP:
		return nil
	case plumb.SYNC:
		return c.Get(p.URL(), p.Path(), p.Ref())
	}
	return fmt.Errorf(`Unrecognized  method %d`, p.Method())
}

// GetTemplate fetch template and store in template store
func (c *Client) GetTemplate(url, ref string) (*plumb.Plumb, error) {
	p := c.plumbTemplates(url, ref)
	return p, c.GetPlumb(p)
}

// GetRepo fetch repo and store in repo store
func (c *Client) GetRepo(url, ref string) (*plumb.Plumb, error) {
	p := c.plumbRepos(url, ref)
	return p, c.GetPlumb(p)
}

// Sync will download/synchronize with the upstream git repo
func (d *Client) Sync(url, path, ref string) error {
	var err error
	err = d.Clone(url, path)
	if err != nil {
		return fmt.Errorf("sync err: %w", err)
	}
	err = d.Pull(path)
	if err != nil {
		return fmt.Errorf("sync err: %w", err)
	}
	err = d.Checkout(path, ref)
	if err != nil {
		return fmt.Errorf("sync err: %w", err)
	}
	return nil
}

// Clone will clone a remote repository
func (d *Client) Clone(url, path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
		})
		if d.err.LogF("Can not clone %s: %v", url, err) {
			return fmt.Errorf("clone error: %w", err)
		}
	}
	return nil
}

// Pull will pull from the remote repository
func (d *Client) Pull(path string) error {
	r, err := git.PlainOpen(path)
	if d.err.LogF("Error opening path '%s': %+v", path, err) {
		return fmt.Errorf("pull error: %w", err)
	}

	w, err := r.Worktree()
	if d.err.LogF("Error reading path '%s': %+v", path, err) {
		return fmt.Errorf("pull error: %w", err)
	}

	pullopts := &git.PullOptions{}
	err = w.Pull(pullopts)
	if err != git.NoErrAlreadyUpToDate {
		if d.err.LogF("Error pulling %s: %+v", path, err) {
			return fmt.Errorf("pull error: %w", err)
		}
	}
	return nil
}

// Tags will list all tags
func (d *Client) Tags(path string) []string {
	taglist := []string{}

	r := d.plainopen(path)
	iter, err := r.Tags()
	d.err.LogF("Error listing tags %s: %+v", path, err)

	fn := func(tag *plumbing.Reference) error {
		taglist = append(taglist, tag.Name().Short())
		return nil
	}
	err = iter.ForEach(fn)
	d.err.LogF("Error listing tags %s: %+v", path, err)

	return taglist
}

// Checkout checks out a reference. If ref is an empty string will checkout using the internally set ref
func (d *Client) Checkout(path, ref string) error {
	if ref == "" {
		return nil
	}
	r, err := git.PlainOpen(path)
	if d.err.LogF("Error opening path '%s': %+v", path, err) {
		return fmt.Errorf("checkout error: %w", err)
	}

	refObj, err := r.Reference(plumbing.ReferenceName(ref), true)
	if d.err.LogF("Error reading reference '%s' for path '%s': %+v", ref, path, err) {
		return fmt.Errorf("checkout error: %w", err)
	}

	chkops := &git.CheckoutOptions{
		Hash: refObj.Hash(),
	}

	w, err := r.Worktree()
	if d.err.LogF("Error reading path '%s': %+v", path, err) {
		return fmt.Errorf("checkout error: %w", err)
	}

	err = w.Checkout(chkops)
	if d.err.LogF("Error checkout out: %+v", err) {
		return fmt.Errorf("checkout error: %w", err)
	}
	return nil
}

func (d *Client) plainopen(path string) *git.Repository {
	if path == "" {
		return nil
	}
	r, err := git.PlainOpen(path)
	d.err.LogF("Error opening path '%s': %+v", path, err)
	return r
}
