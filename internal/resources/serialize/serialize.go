package serialize

// RepoMain yaml file stored as `repo.yml` in the projects root directory
type RepoMain struct {
	Name         string   `yaml:"name" validate:"required,alphanum"`
	Desc         string   `yaml:"description" validate:"required"`
	TemplateURLs []string `yaml:"templates"`
}

// RepoTemplateFile file written to a repo as `template/${TEMPLATE}.yml`
type RepoTemplateFile struct {
	Name     string   `yaml:"name" validate:"required,alphanum"`
	Desc     string   `yaml:"description" validate:"required"`
	URL      string   `yaml:"url" validate:"required,url"`
	Versions []string `yaml:"versions"`
}
