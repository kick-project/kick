package client

import (
	"fmt"
	"io"

	"github.com/go-playground/validator"
	"github.com/kick-project/kick/internal/di/callbacks"
	"github.com/kick-project/kick/internal/resources/client/plumb"
	"github.com/kick-project/kick/internal/resources/errs"
	"github.com/kick-project/kick/internal/resources/vcs"
)

// Client
type Client struct {
	err            errs.HandlerIface
	plumbRepos     callbacks.MakePlumb
	plumbTemplates callbacks.MakePlumb
	stdout         io.Writer
	vcs            *vcs.VCS
}

// Options for New function
type Options struct {
	CallPlumbRepos     callbacks.MakePlumb `validate:"required"`
	CallPlumbTemplates callbacks.MakePlumb `validate:"required"`
	Err                errs.HandlerIface   `validate:"required"`
	Stdout             io.Writer           `validate:"required"`
	VCS                *vcs.VCS            `validate:"required"`
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
		plumbRepos:     opts.CallPlumbRepos,
		plumbTemplates: opts.CallPlumbTemplates,
		stdout:         opts.Stdout,
		vcs:            opts.VCS,
	}
}

// Get get url and clone/sync to path using ref.
// Defaults to default branch if ref is nil.
func (c *Client) Get(url, path, ref string) error {
	repo, err := c.vcs.Clone(url, path)
	if err != nil {
		return fmt.Errorf("get clone error: %w", err)
	}
	err = repo.Checkout(ref)
	if err != nil {
		return fmt.Errorf("get checkout error: %w", err)
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
	p, err := c.plumbTemplates(url, ref)
	if err != nil {
		return nil, err
	}
	return p, c.GetPlumb(p)
}

// GetRepo fetch repo and store in repo store
func (c *Client) GetRepo(url, ref string) (*plumb.Plumb, error) {
	p, err := c.plumbRepos(url, ref)
	if err != nil {
		return nil, err
	}
	return p, c.GetPlumb(p)
}
