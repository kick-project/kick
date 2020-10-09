package entry

// Entry an individual entry for a given search
type Entry struct {
	Name       string // Template name
	URL        string // URL location
	Desc       string // Description
	MasterName string // The master associated with the template
	MasterURL  string // The masters' URL location
	MasterDesc string // The masters' description
}
