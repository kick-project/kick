package iplumbing

import "github.com/kick-project/kick/internal/settings"

// InjectGlobal injects settings for plumbing.Plumb
func InjectGlobal(s *settings.Settings) (opts struct {
	Base string
}) {
	opts.Base = s.PathGlobalDir
	return
}

// InjectMaster injects settings for plumbing.Plumb
func InjectMaster(s *settings.Settings) (opts struct {
	Base string
}) {
	opts.Base = s.PathMasterDir
	return
}

// InjectTemplate injects settings for plumbing.Plumb
func InjectTemplate(s *settings.Settings) (opts struct {
	Base string
}) {
	opts.Base = s.PathTemplateDir
	return
}
