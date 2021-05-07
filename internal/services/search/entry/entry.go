package entry

// Entry an individual entry for a given search
type Entry struct {
	Name       string // Template name
	URL        string // URL location
	Desc       string // Description
	RepoName string // The repo associated with the template
	RepoURL  string // The repos' URL location
	RepoDesc string // The repos' description
}
