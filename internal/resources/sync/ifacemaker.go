// DO NOT EDIT: Generated using "make interfaces"

package sync

// SyncIface ...
type SyncIface interface {
	// Repo syncs repo data
	Repo()
	// Files synchronizes templates between the YAML configuration, database
	// and its upstream version control repository.
	Files()
}
