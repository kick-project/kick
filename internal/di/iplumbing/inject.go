package iplumbing

import "github.com/kick-project/kick/internal/di"

// InjectMaster injects di for plumbing.Plumb
func InjectMaster(s *di.DI) (opts struct {
	Base string
}) {
	opts.Base = s.PathMasterDir
	return
}

// InjectTemplate injects di for plumbing.Plumb
func InjectTemplate(s *di.DI) (opts struct {
	Base string
}) {
	opts.Base = s.PathTemplateDir
	return
}
