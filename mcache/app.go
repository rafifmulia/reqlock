package mcache

import (
	"reflect"
	"sync"
	"time"
)

var (
	data cacheKey    = make(cacheKey)
	gomu *sync.Mutex = &sync.Mutex{}
)

type cacheKey map[string]*cacheValue
type cacheValue struct {
	v   map[any]bool // value of cache that turns into hash.
	mtm int64        // time.Now().Unix()
	mu  *sync.RWMutex
}

// Initialize [data] if its nil.
func initCacheKey(k string) {
	gomu.Lock()
	defer gomu.Unlock()
	if data[k] == nil {
		setDfCacheKey(k)
	}
}

// Set default value of [data].
func setDfCacheKey(k string) {
	data[k] = &cacheValue{
		v:  make(map[any]bool),
		mu: &sync.RWMutex{},
	}
}

// Set cache, and return false if cache is already exist.
// k is cacheKey, and v is your cacheValue.
func Set(k string, v any) bool {
	var (
		pv reflect.Value
	)
	initCacheKey(k)
	pv = reflect.ValueOf(v)
	if pv.Kind() == reflect.Ptr {
		v = pv.Elem().Interface() // Dereference if its pointer.
	}
	data[k].mu.Lock()
	defer data[k].mu.Unlock()
	if data[k].v[v] {
		return false
	}
	data[k].v[v] = true
	data[k].mtm = time.Now().Unix()
	return true
}

// Check if cache is exist, return true if its exist.
// k is cacheKey, and v is your cacheValue.
func Check(k string, v any) bool {
	if data[k] == nil {
		return false
	}
	data[k].mu.RLock()
	defer data[k].mu.RUnlock()
	if data[k].v == nil {
		return false
	}
	if !data[k].v[v] {
		return false
	}
	return true
}
