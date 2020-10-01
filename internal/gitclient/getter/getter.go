package getter

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

// Getter Manages downloading of metadata
type Getter struct {
	base   string
	Scheme string
	Path   string
	Branch string
	URL    string
	Method int
}

// New create a new Getter. All repositories downloaded will be under the basedir.
func New(basedir string) *Getter {
	return &Getter{
		base: basedir,
	}
}

// Handler Set the item to get
func (m *Getter) Handler(url string) error {
	urlexp := utils.ExpandPath(url)

	// Local filesystem path do nothing
	if strings.HasPrefix(urlexp, "/") {
		m.Path = urlexp
		m.Method = NOOP
		return nil
	}

	u := &utils.URLx{}
	err := u.Parse(url)
	if err != nil {
		return err
	}
	m.URL = url
	m.Scheme = u.Scheme
	if u.Scheme == "file" {
		m.Path = u.Path
		m.Method = NOOP
		return nil
	}
	m.Path = m.localPath(u)
	m.Method = SYNC

	return nil
}

// localPath determines local path to template
func (m *Getter) localPath(u *utils.URLx) string {
	if u.Scheme == "file" {
		return u.Path
	}
	p := filepath.Join(m.base, u.Path, u.Project)
	return p
}
