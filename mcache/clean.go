package mcache

import (
	"time"
)

var (
	ticker *time.Ticker
)

// Remove specific v from [cacheVal].
func Remove(k string, v any) error {
	var (
		cv *cacheVal
	)
	if data[k] == nil {
		return nil
	}
	cv = data[k]
	cv.mu.Lock()
	defer cv.mu.Unlock()
	delete(data[k].v, v)
	return nil
}

// Clear specific [cacheKey].
func Delete(k string) error {
	gomu.Lock()
	defer gomu.Unlock()
	delete(data, k)
	setDfCacheKey(k)
	return nil
}

// Flush all caches.
func Flush() error {
	var (
		cvs []*cacheVal
	)
	gomu.Lock()
	defer gomu.Unlock()
	for _, cv := range data {
		cvs = append(cvs, cv)
	}
	for i := range cvs {
		cvs[i].mu.Lock()
		defer cvs[i].mu.Unlock()
	}
	clear(data)
	data = make(cacheKey)
	return nil
}

// CleanupRoutine will run every givenTime,
// and remove the item of [cacheVal] if more than n second since its modified.
func CleanupRoutine(givenTime time.Duration, n int64) {
	ticker = time.NewTicker(givenTime)
	for t := range ticker.C {
		for k, v := range data {
			if v == nil {
				continue
			}
			for vv, mtm := range v.v {
				diff := t.Unix() - mtm
				if diff > n {
					Remove(k, vv)
				}
			}
		}
	}
}

func ShutdownCleanupRoutine() error {
	ticker.Stop()
	return nil
}
