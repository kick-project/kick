package iplumbing

import "github.com/crosseyed/prjstart/internal/settings"

// Inject injects settings for plumbing.Plumb
func Inject(s *settings.Settings) (opts struct {
	Base string
}) {
	opts.Base = s.PathTemplateDir
	return
}
