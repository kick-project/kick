package db

import (
	"sync"
)

// mu Mutex used to co-ordinate all writes to a sequential process
// instead of parallel as per Sqlite3 documentation.
var mu = &sync.Mutex{}

// Lock locks using a shared sync.Mutex{} in the db package.
func Lock() {
	mu.Lock()
}

// Unlock unlocks using a shared sync.Mutex{} in the db package.
func Unlock() {
	mu.Unlock()
}
