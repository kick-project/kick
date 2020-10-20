package ilist

import (
	"io"

	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/settings"
)

// Inject create settings for list.List
func Inject(s *settings.Settings) (opts struct {
	Stderr io.Writer
	Stdout io.Writer
	Conf   *config.File
}) {
	opts.Stderr = s.Stderr
	opts.Stdout = s.Stdout
	opts.Conf = s.ConfigFile()
	return
}
