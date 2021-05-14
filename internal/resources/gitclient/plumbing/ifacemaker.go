// DO NOT EDIT: Generated using "make interfaces"

package plumbing

// PlumbingIface ...
type PlumbingIface interface {
	// Handler Set the item to get
	Handler(url string) error
	// Scheme URL scheme.
	Scheme() string
	// Path local path on disk.
	Path() string
	// Branch branch to checkout.
	Branch() string
	// URL original URL.
	URL() string
	// Method actions to perform.
	Method() int
	// Local takes relative path returns absolute path.
	// Slash is replaced using path seperator.
	Local(relative string) string
}
