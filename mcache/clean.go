package mcache

import "time"

var (
	ticker *time.Ticker
)

// Clear specific cache.
func Delete(k string) error {
	mu.Lock()
	defer mu.Unlock()
	data[k] = new(objKey)
	return nil
}

// Flush all caches.
func Flush() error {
	mu.Lock()
	defer mu.Unlock()
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
			mu.RLock()
			diff := t.Unix() - v.mtm
			mu.RUnlock()
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
