package serialize

// RepoMain yaml file stored as `repo.yml` in the projects root directory
type RepoMain struct {
	Name         string   `yaml:"name" validate:"required,alphanum"`
	Desc         string   `yaml:"description" validate:"required"`
	TemplateURLs []string `yaml:"templates"`
}

// TemplateMain template yaml file stored as `.kick.yml` in the projects root directory
type TemplateMain struct {
	Name string `yaml:"name" validate:"required,alphanum"`
	Desc string `yaml:"description" validate:"required"`
}

// RepoTemplateFile file written to a repo as `template/${TEMPLATE}.yml`
type RepoTemplateFile struct {
	Name string `yaml:"name" validate:"required,alphanum"`
	Desc string `yaml:"description" validate:"required"`
	URL  string `yaml:"url" validate:"required,url"`
}
