package globals

import (
	"github.com/crosseyed/prjstart/internal/resources/config"
	"github.com/crosseyed/prjstart/internal/services/template"
)

// Config serves as the main configuration piece
var Config *config.File

// Vars all template variables
var Vars *template.Variables
