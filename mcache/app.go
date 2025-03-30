package mcache

import "sync"

var (
	data map[string]map[any]bool = make(map[string]map[any]bool)
	mux  *sync.Mutex             = &sync.Mutex{}
)

// Return false if key and value is already exist.
func Set(k string, v any) bool {
	mux.Lock()
	defer mux.Unlock()
	if data[k][v] {
		return false
	}
	data[k] = make(map[any]bool)
	data[k][v] = true
	return true
}
