// DO NOT EDIT: Generated using "make interfaces"

package gitclient

// GitclientIface ...
type GitclientIface interface {
	// Sync will download/synchronize with the upstream git repo
	Sync()
	// SetRef sets the default reference. See Checkout
	SetRef(ref string)
	// Clone will clone a remote repository
	Clone()
	// Pull will pull from the remote repository
	Pull()
	// Tags will list all tags
	Tags() []string
	// Checkout checks out a reference. If ref is an empty string will checkout using the internally set ref
	Checkout(ref string)
}
