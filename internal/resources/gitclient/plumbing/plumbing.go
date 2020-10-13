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
	Base string
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
		Base: basedir,
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
	if plu.mb == nil {
		return ""
	}
	return plu.mb.scheme
}

// Path local path on disk.
func (plu *Plumbing) Path() string {
	if plu.mb == nil {
		return ""
	}
	return plu.mb.path
}

// Branch branch to checkout.
func (plu *Plumbing) Branch() string {
	if plu.mb == nil {
		return ""
	}
	return plu.mb.branch
}

// URL original URL.
func (plu *Plumbing) URL() string {
	if plu.mb == nil {
		return ""
	}
	return plu.mb.uRL
}

// Method actions to perform.
func (plu *Plumbing) Method() int {
	if plu.mb == nil {
		return 0
	}
	return plu.mb.method
}

// localPath determines local path to template
func (plu *Plumbing) localPath(u *utils.URLx) string {
	if u.Scheme == "file" {
		return u.Path
	}
	lPath := filepath.Join(plu.Base, u.Path, u.Project)
	return lPath
}
