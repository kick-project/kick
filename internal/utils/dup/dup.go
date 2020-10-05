package dup

const (
	// OK Insert ok
	OK = iota + 1
	// DUP Duplicate
	DUP
)

// String checks for duplicate strings
type String struct {
	list map[string]bool
}

// Check Inserts a string and checks if one has be previously inserted
func (ds *String) Check(item string) int {
	if ds.list == nil {
		ds.list = map[string]bool{}
	}
	if _, dup := ds.list[item]; dup {
		return DUP
	}

	ds.list[item] = true
	return OK
}
