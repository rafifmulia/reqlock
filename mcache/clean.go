package mcache

import (
	"time"
)

var (
	ticker *time.Ticker
)

// Clear specific cache.
func Delete(k string) error {
	gomu.Lock()
	defer gomu.Unlock()
	delete(data, k)
	setDfCacheKey(k)
	return nil
}

// Flush all caches.
func Flush() error {
	gomu.Lock()
	defer gomu.Unlock()
	// Try to prevent `defer data[k].mu.Unlock()` from nil pointer dereference in [Set].
	for _, cv := range data {
		cv.mu.Lock()
		defer cv.mu.Unlock()
	}
	clear(data)
	data = make(cacheKey)
	return nil
}

// CleanupRoutine will run every givenTime,
// and clean cacheKey if not used since n second.
func CleanupRoutine(givenTime time.Duration, n int64) {
	ticker = time.NewTicker(givenTime)
	for t := range ticker.C {
		for k, v := range data {
			if v == nil {
				continue
			}
			v.mu.RLock()
			diff := t.Unix() - v.mtm
			v.mu.RUnlock()
			if diff > n {
				Delete(k)
			}
		}
	}
}

func ShutdownCleanupRoutine() error {
	ticker.Stop()
	return nil
}
