package internal

type Installer struct {
	set      string
	template string
	name     string
}

// getSet returns the cached set path
func (s *Installer) getSet() (string, error) {
	return "", nil
}

// addTemplate adds the template defined in set if it exists
func (s *Installer) addTemplate() error {
	return nil
}

//
// Errors
//

type ErrorNoset struct {
	error
}

func Error() string {
	return "no set found"
}
