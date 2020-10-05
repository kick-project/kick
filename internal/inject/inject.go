package inject

import "sync"

var once = sync.Once{}
var mu = sync.Mutex{}

// store is where all injected dependencies
var store map[string]interface{} = map[string]interface{}{}

// Set sets a interface
func Set(handle string, item interface{}) {
	Lock()
	defer Unlock()
	store[handle] = item
}

// Get returns an injected object
func Get(handle string) interface{} {
	item, ok := store[handle]
	if ok {
		return item
	}
	return nil
}

// Init initializes all dependencies
func Init() {
	once.Do(_init)
}

// Reset resets stored objects
func Reset() {
	_init()
}

// Lock lock out all changes to object state.
func Lock() {
	mu.Lock()
}

// Unlock unlock all changes to object state.
func Unlock() {
	mu.Unlock()
}

func _init() {
	Lock()
	defer Unlock()
	store = map[string]interface{}{}
}
