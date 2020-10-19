package iplumbing

import "github.com/kick-project/kick/internal/settings"

// Inject injects settings for plumbing.Plumb
func Inject(s *settings.Settings) (opts struct {
	Base string
}) {
	opts.Base = s.PathTemplateDir
	return
}
