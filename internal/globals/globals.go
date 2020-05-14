package globals

import (
	"github.com/crosseyed/prjstart/internal/config"
	"github.com/crosseyed/prjstart/internal/template"
)

// Config serves as the main configuration piece
var Config *config.Config

// Vars all template variables
var Vars *template.TmplVars
