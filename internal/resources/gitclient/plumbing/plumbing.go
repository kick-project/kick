package plumbing

import (
	"path/filepath"
	"strings"

	"github.com/crosseyed/prjstart/internal/utils"
)

const (
	// NOOP do not do anything
	NOOP = iota + 1
	// SYNC Git synchronize (Clone, Pull, Checkout)
	SYNC
)

// Plumbing Manages downloading of metadata
type Plumbing struct {
	base string
	mb   *plumb
}

type plumb struct {
	scheme string
	path   string
	branch string
	uRL    string
	method int
}

// New create a new plumber. All repositories downloaded will be under the basedir.
func New(basedir string) *Plumbing {
	return &Plumbing{
		base: basedir,
		mb:   &plumb{},
	}
}

// Handler Set the item to get
func (plu *Plumbing) Handler(url string) error {
	urlexp := utils.ExpandPath(url)
	gt := &plumb{}

	// Local filesystem path do nothing
	if strings.HasPrefix(urlexp, "/") {
		gt.path = urlexp
		gt.method = NOOP
		plu.mb = gt
		return nil
	}

	u := &utils.URLx{}
	err := u.Parse(url)
	if err != nil {
		plu.mb = &plumb{}
		return err
	}
	gt.uRL = url
	gt.scheme = u.Scheme
	if u.Scheme == "file" {
		gt.path = u.Path
		gt.method = NOOP
		plu.mb = gt
		return nil
	}
	gt.path = plu.localPath(u)
	gt.method = SYNC
	plu.mb = gt

	return nil
}

// Scheme URL scheme.
func (plu *Plumbing) Scheme() string {
	return plu.mb.scheme
}

// Path local path on disk.
func (plu *Plumbing) Path() string {
	return plu.mb.path
}

// Branch branch to checkout.
func (plu *Plumbing) Branch() string {
	return plu.mb.branch
}

// URL original URL.
func (plu *Plumbing) URL() string {
	return plu.mb.uRL
}

// Method actions to perform.
func (plu *Plumbing) Method() int {
	return plu.mb.method
}

// localPath determines local path to template
func (plu *Plumbing) localPath(u *utils.URLx) string {
	if u.Scheme == "file" {
		return u.Path
	}
	lPath := filepath.Join(plu.base, u.Path, u.Project)
	return lPath
}
