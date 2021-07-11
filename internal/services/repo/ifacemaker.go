// DO NOT EDIT: Generated using "make interfaces"

package repo

// RepoIface ...
type RepoIface interface {
	// Build build repo
	Build()
	// List list repositories
	List()
	// Info information on repositories
	Info(repo string)
}
