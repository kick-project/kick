package cond

// ContainsString returns true if any of hte strings match
func ContainsString(s string, strs ...string) bool {
	for _, x := range strs {
		if x == s {
			return true
		}
	}
	return false
}
