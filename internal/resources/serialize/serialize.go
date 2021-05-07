package serialize

// Repo yaml file
type Repo struct {
	Name         string   `yaml:"name"`
	Desc         string   `yaml:"description"`
	TemplateURLs []string `yaml:"templates"`
}

// Kick template yaml file
type Kick struct {
	Name string `yaml:"name"`
	Desc string `yaml:"description"`
}

// RepoElement file
type RepoElement struct {
	Name string `yaml:"name"`
	Desc string `yaml:"description"`
	URL  string `yaml:"url"`
}

// TemplateElement file
type TemplateElement struct {
	Name string `yaml:"name"`
	Desc string `yaml:"description"`
	URL  string `yaml:"url"`
}
