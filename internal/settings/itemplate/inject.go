package itemplate

import (
	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/services/template/variables"
	"github.com/crosseyed/prjstart/internal/settings"
)

// Inject creates settings for template.New
func Inject(s *settings.Settings) (opts struct {
	Config      *config.File
	Variables   *variables.Variables
	TemplateDir string
	ModeLineLen uint8
}) {
	configFile := s.ConfigFile()
	vars := variables.New()
	vars.ProjectVariable("name", s.ProjectName)
	opts.Config = configFile
	opts.TemplateDir = s.PathTemplateDir
	opts.Variables = vars

	return opts
}
