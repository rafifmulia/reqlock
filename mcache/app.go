package mcache

import (
	"reflect"
	"sync"
	"time"
)

var (
	data cacheKey      = make(cacheKey)
	mu   *sync.RWMutex = &sync.RWMutex{}
)

type cacheKey map[string]*objKey
type objKey struct {
	mtm int64 // time.Now().Unix()
	v   map[any]bool
}

// v must not pointer. If v is pointer, then the value will be copied.
// Return false if key and value is already exist.
func Set(k string, v any) bool {
	var (
		vo reflect.Value
	)
	vo = reflect.ValueOf(v)
	if vo.Kind() == reflect.Ptr {
		v = vo.Elem().Interface()
	}
	mu.Lock()
	defer mu.Unlock()
	if data[k] == nil {
		data[k] = new(objKey)
	}
	if data[k].v == nil {
		data[k].v = make(map[any]bool)
	}
	if data[k].v[v] {
		return false
	}
	data[k].v[v] = true
	data[k].mtm = time.Now().Unix()
	return true
}

// Return true if its exist.
func Check(k string, v any) bool {
	mu.RLock()
	defer mu.RUnlock()
	if data[k] == nil {
		return false
	}
	if data[k].v == nil {
		return false
	}
	if !data[k].v[v] {
		return false
	}
	return true
}
