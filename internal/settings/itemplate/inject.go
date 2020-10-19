package itemplate

import (
	"io"

	"github.com/apex/log"
	"github.com/kick-project/kick/internal/resources/config"
	"github.com/kick-project/kick/internal/resources/template/renderer"
	"github.com/kick-project/kick/internal/resources/template/variables"
	"github.com/kick-project/kick/internal/settings"
)

// Inject creates settings for template.New
func Inject(s *settings.Settings) (opts struct {
	Config         *config.File
	Log            *log.Logger
	ModeLineLen    uint8
	RenderersAvail map[string]renderer.Renderer
	RenderCurrent  string
	Stderr         io.Writer
	Stdout         io.Writer
	TemplateDir    string
	Variables      *variables.Variables
}) {
	configFile := s.ConfigFile()
	vars := variables.New()
	vars.ProjectVariable("NAME", s.ProjectName)
	opts.Config = configFile
	opts.Log = s.GetLogger()
	opts.Stderr = s.Stderr
	opts.Stdout = s.Stdout
	opts.TemplateDir = s.PathTemplateDir
	opts.Variables = vars
	opts.RenderCurrent = "envsubst"
	opts.RenderersAvail = map[string]renderer.Renderer{
		"texttemplate": &renderer.RenderText{},
		"envsubst":     &renderer.RenderEnv{},
	}

	return opts
}
