// Package fflags enables featureflags using environment variables. Feature
// flags can be globally disabled/enabled using the FF_ENABLED environment
// variable (true enables any other value disables).
package fflags

import "os"

// enabled toggle to determine if feature flags have been turned on.
// This is needed so if anyone accidentally sets one variable the feature will still
// not be turned on. It also makes it easier to toggle all feature flags.
func enabled() bool {
	return os.Getenv("FF_ENABLED") == "true"
}

// Remote enables remote operations
func Remote() bool {
	return enabled() && os.Getenv("FF_KICK_REMOTE") == "true"
}

// GitClone refactor the way git repositories are cloned
func GitClone() bool {
	return enabled() && os.Getenv("FF_KICK_GIT_CLONE") == "true"
}

// ORM create ORM based schema instead of through SQL
func ORM() bool {
	return enabled() && os.Getenv("FF_KICK_ORM") == "true"
}
