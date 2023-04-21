package handle

import (
	"errors"

	"github.com/kick-project/kick/internal/di/callbacks"
	"github.com/kick-project/kick/internal/resources/config"
)

type Handle struct {
	config config.File
	plumb  callbacks.MakePlumb
}

type Options struct {
	Config config.File
	Plumb  callbacks.MakePlumb
}

func New(opts Options) *Handle {
	h := &Handle{
		config: opts.Config,
		plumb:  opts.Plumb,
	}
	return h
}

func (h *Handle) Handle2Template(handle string) *config.Template {
	for _, tConf := range h.config.Templates {
		if tConf.Handle == handle {
			return &tConf
		}
	}
	return nil
}

var ErrNoHandle = errors.New("handle not found")

func (h *Handle) Handle2Path(handle string) (string, error) {
	t := h.Handle2Template(handle)
	if t == nil {
		return "", ErrNoHandle
	}
	p, err := h.plumb(t.URL, "")
	if err != nil {
		return "", err
	}
	return p.Path(), nil
}
