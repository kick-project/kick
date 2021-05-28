// Package env holds environment variables used by the application. They differ
// in use to the fflags package which also uses environment variables, but is
// only to enable or disable feature flags.
package env

import "os"

type Vars struct {
}

//
// Options
//

func (v *Vars) LogFile() string {
	return os.Getenv("KICK_LOG")
}

//
// Development
//

// Debug Turn debug logging on. See di.DI
func (v *Vars) Debug() bool {
	return os.Getenv("KICK_DEBUG") == "true"
}
