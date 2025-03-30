package mcache

import (
	"reflect"
	"sync"
)

var (
	data map[string]map[any]bool = make(map[string]map[any]bool)
	mux  *sync.Mutex             = &sync.Mutex{}
)

// v must not pointer.
// Return false if key and value is already exist.
func Set(k string, v any) bool {
	var (
		vo reflect.Value
	)
	vo = reflect.ValueOf(v)
	if vo.Kind() == reflect.Ptr {
		v = vo.Elem().Interface()
	}
	mux.Lock()
	defer mux.Unlock()
	if data[k][v] {
		return false
	}
	if data[k] == nil {
		data[k] = make(map[any]bool)
	}
	data[k][v] = true
	return true
}
