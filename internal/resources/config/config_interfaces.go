// AUTO GENERATED. DO NOT EDIT

package config

// FileIface ...
type FileIface interface {
	// AppendTemplate appends a template to list of templates.
	// If stop is non zero, the calling function should exit the program with the
	// value contained in stop.
	AppendTemplate(t Template) (err error)
	// Load loads configuration file from disk
	Load() error
	// SaveTemplates saves template configuration file to disk
	SaveTemplates() error
}
