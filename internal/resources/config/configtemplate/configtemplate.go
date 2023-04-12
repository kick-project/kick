package configtemplate

// TemplateMain template yaml file stored as `.kick.yml` in the projects root directory
type TemplateMain struct {
	Name   string              `yaml:"name" validate:"required,alphanum"`
	Desc   string              `yaml:"description" validate:"required"`
	Envs   map[string]string   `yaml:"envs"` // Required environment variables
	Labels map[string][]string `yaml:"label"`
}
