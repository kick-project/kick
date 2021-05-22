package plumb

import (
	"path/filepath"
	"strings"

	"github.com/kick-project/kick/internal/resources/file"
	"github.com/kick-project/kick/internal/resources/parse"
)

const (
	// NOOP do not do anything
	NOOP = iota + 1
	// SYNC Git synchronize (Clone, Pull, Checkout)
	SYNC
)

// Plumb plumbing for fetching URLs
type Plumb struct {
	base   string
	url    string
	scheme string
	path   string
	ref    string
	method int
}

// New is a Plumb constructor
func New(basedir, url, ref string) *Plumb {
	p := &Plumb{
		base: basedir,
		url:  url,
		ref:  ref,
	}
	err := p.parse(url)
	if err != nil {
		panic(err)
	}
	return p
}

// parse parse url
func (p *Plumb) parse(url string) error {
	urlexp := file.ExpandPath(url)

	// Local filesystem path do nothing
	if strings.HasPrefix(urlexp, "/") {
		p.path = urlexp
		p.method = NOOP
		return nil
	}

	u := &parse.URLx{}
	err := u.Parse(url)
	if err != nil {
		return err
	}
	p.url = url
	p.scheme = u.Scheme
	if u.Scheme == "file" {
		p.path = u.Path
		p.method = NOOP
		return nil
	}
	p.path = p.localPath(u)
	p.method = SYNC

	return nil
}

// Scheme URL scheme.
func (p *Plumb) Scheme() string {
	return p.scheme
}

// Path local path on disk.
func (p *Plumb) Path() string {
	return p.path
}

// Ref branch to checkout.
func (p *Plumb) Ref() string {
	return p.ref
}

// URL original URL.
func (p *Plumb) URL() string {
	return p.url
}

// Method actions to perform.
func (p *Plumb) Method() int {
	return p.method
}

// Local takes relative path returns absolute path.
// Slash is replaced using path seperator.
func (p *Plumb) Local(relative string) string {
	return filepath.Join(p.Path(), filepath.FromSlash(relative))
}

// localPath determines local path to template
func (p *Plumb) localPath(u *parse.URLx) string {
	if u.Scheme == "file" {
		return u.Path
	}
	lPath := filepath.Join(p.base, u.Path, u.Project)
	return lPath
}
