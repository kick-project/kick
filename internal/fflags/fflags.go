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
	return enabled() && os.Getenv("FF_PRJSTART_REMOTE") == "true"
}

// DB mocks the database
func DB() bool {
	return enabled() && os.Getenv("FF_PRJSTART_MOCK_DB") == "true"
}

// GitClone refactor the way git repositories are cloned
func GitClone() bool {
	return enabled() && os.Getenv("FF_PRJSTART_GIT_CLONE") == "true"
}
