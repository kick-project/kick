// Package env holds environment variables used by the application. They differ
// in use to the fflags package which also uses environment variables, but is
// only to enable or disable feature flags.
package env

import "os"

// Debug Turn debug logging on. See settings.Settings
func Debug() bool {
	return os.Getenv("PRJSTART_DEBUG") == "true"
}
