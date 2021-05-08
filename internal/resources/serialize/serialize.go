package serialize

// RepoMain yaml file stored as `repo.yml` in the projects root directory
type RepoMain struct {
	Name         string   `yaml:"name"`
	Desc         string   `yaml:"description"`
	TemplateURLs []string `yaml:"templates"`
}

// TemplateMain template yaml file stored as `.kick.yml` in the projects root directory
type TemplateMain struct {
	Name string `yaml:"name"`
	Desc string `yaml:"description"`
}

// Repo file
type Repo struct {
	Name string `yaml:"name"`
	Desc string `yaml:"description"`
	URL  string `yaml:"url"`
}

// Template file
type Template struct {
	Name string `yaml:"name"`
	Desc string `yaml:"description"`
	URL  string `yaml:"url"`
}
