package iplumbing

import "github.com/kick-project/kick/internal/di"

// InjectRepo injects di for plumbing.Plumb
func InjectRepo(s *di.DI) (opts struct {
	Base string
}) {
	opts.Base = s.PathRepoDir
	return
}

// InjectTemplate injects di for plumbing.Plumb
func InjectTemplate(s *di.DI) (opts struct {
	Base string
}) {
	opts.Base = s.PathTemplateDir
	return
}
