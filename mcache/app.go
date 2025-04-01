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

type cacheKey map[string]*cacheVal
type cacheVal struct {
	v  map[any]int64 // map key is a value of cache that turns into hash, and int64 is time.Now().Unix().
	mu *sync.RWMutex
}

// Initialize [data] if its nil.
func initCacheKey(k string) {
	if data[k] == nil {
		setDfCacheKey(k)
	}
}

// Set default value of [data].
func setDfCacheKey(k string) {
	data[k] = &cacheVal{
		v:  make(map[any]int64),
		mu: &sync.RWMutex{},
	}
}

// Set cache, and return false if cache is already exist.
// k is cacheKey, and v is your cacheVal.
func Set(k string, v any) bool {
	var (
		pv reflect.Value
		cv *cacheVal
	)
	gomu.Lock()
	initCacheKey(k)
	cv = data[k]
	gomu.Unlock()
	pv = reflect.ValueOf(v)
	if pv.Kind() == reflect.Ptr {
		v = pv.Elem().Interface() // Dereference if its pointer.
	}
	cv.mu.Lock()
	defer cv.mu.Unlock()
	if cv.v[v] > 0 {
		return false
	}
	cv.v[v] = time.Now().Unix()
	return true
}

// Check if cache is exist, return true if its exist.
// k is cacheKey, and v is your cacheVal.
func Check(k string, v any) bool {
	if data[k] == nil {
		return false
	}
	data[k].mu.RLock()
	defer data[k].mu.RUnlock()
	if data[k].v == nil {
		return false
	}
	if data[k].v[v] < 1 {
		return false
	}
	return true
}
