package mcache

import (
	"sync"
	"testing"
	"time"
)

// Internal types, variables, and functions for testing only.
type ticket struct {
	id, room, seat int32
	film           string
}

var (
	bookseat string = "bookseat"                               // cacheKey
	dfTicket ticket = ticket{film: "batman", room: 2, seat: 7} // Default ticket.
)

// checkCache will check either duplicated caches or unsuccessful to set caches (even single cache).
// statuses is return values from [Set].
// Duplicated caches happens when total caches that has been succeded to [Set] (returns true),
// bigger than maxCache.
// Unsuccessful to set caches happens when count is 0.
func checkCache(statuses []bool, maxCache uint32) {
	var count uint32
	for _, v := range statuses {
		if v {
			count++
		}
	}
	if count > maxCache {
		panic("found duplicated set")
	}
	if count < maxCache {
		panic("no cache has been set")
	}
}

// go test -v -count=1 -failfast -cpu=1 -run='^TestDefaultValue$'
func TestDefaultValue(t *testing.T) {
	k := bookseat
	println(data)    // 0xc0000a2300
	println(data[k]) // 0x0 == nil
	data[k] = new(cacheValue)
	println(data[k].mu) // 0x0 == nil
	data[k].mu = &sync.RWMutex{}
	println(data[k].mtm) // 0
	data[k].mtm = time.Now().Unix()
	println(data[k].v) // 0x0 == nil
	v := ticket{}
	data[k].v = make(map[any]bool)
	data[k].v[v] = true
}

// [Set] with value as pointer.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -race -run='^TestSet1$'
func TestSet1(t *testing.T) {
	var (
		concurrentReq int32           = 40
		wg            *sync.WaitGroup = &sync.WaitGroup{}
		statuses      []bool          = make([]bool, concurrentReq)
		maxCache      uint32          = 1 // Should only 1 cache that has been succeded to Set (returns true).
	)
	wg.Add(int(concurrentReq))
	for i := int32(0); i < concurrentReq; i++ {
		go func() {
			defer wg.Done()
			statuses = append(statuses, Set(bookseat, &ticket{film: "batman", room: 2, seat: 7}))
		}()
	}
	wg.Wait()
	checkCache(statuses, maxCache)
}

// [Set] with value is pass by value.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -race -run='^TestSet2$'
func TestSet2(t *testing.T) {
	var (
		concurrentReq int32           = 10
		wg            *sync.WaitGroup = &sync.WaitGroup{}
		statuses      []bool          = make([]bool, concurrentReq)
		maxCache      uint32          = 1 // Should only 1 cache that has been succeded to Set (returns true).
	)
	wg.Add(int(concurrentReq))
	for i := int32(0); i < concurrentReq; i++ {
		go func() {
			defer wg.Done()
			statuses = append(statuses, Set(bookseat, ticket{film: "batman", room: 2, seat: 7}))
		}()
	}
	wg.Wait()
	checkCache(statuses, maxCache)
}

// [Set] with value as pointer and another value is pass by value.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -race -run='^TestSet3$'
func TestSet3(t *testing.T) {
	var (
		concurrentReq int32           = 40
		wg            *sync.WaitGroup = &sync.WaitGroup{}
		statuses      []bool          = make([]bool, concurrentReq)
		maxCache      uint32          = 1 // Should only 1 cache that has been succeded to Set (returns true).
	)
	wg.Add(int(concurrentReq + 1))
	go func() {
		defer wg.Done()
		statuses = append(statuses, Set(bookseat, &ticket{film: "batman", room: 2, seat: 7}))
	}()
	for i := int32(0); i < concurrentReq; i++ {
		go func() {
			defer wg.Done()
			statuses = append(statuses, Set(bookseat, ticket{film: "batman", room: 2, seat: 7}))
		}()
	}
	wg.Wait()
	checkCache(statuses, maxCache)
}

// Simple scenario for flush test case.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -run='^TestFlush1$'
func TestFlush1(t *testing.T) {
	t.Run("TestSet3", TestSet3)
	isAvail := func() string {
		if Check(bookseat, dfTicket) {
			return "avail"
		} else {
			return "unavailable"
		}
	}
	if s := isAvail(); s != "avail" {
		t.Fatalf("expected avail, got %s\n", s)
	}
	Flush()
	if s := isAvail(); s != "unavailable" {
		t.Fatalf("expected unavailable, got %s\n", s)
	}
}

// Being flushed when there are another [Set].
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -run='^TestFlush2$'
func TestFlush2(t *testing.T) {
	var (
		concurrentReq int32           = 1000
		wg            *sync.WaitGroup = &sync.WaitGroup{}
	)
	isAvail := func() string {
		if Check(bookseat, dfTicket) {
			return "avail"
		} else {
			return "unavailable"
		}
	}
	wg.Add(int(concurrentReq))
	for i := int32(0); i < concurrentReq; i++ {
		go func() {
			defer wg.Done()
			Set(bookseat, dfTicket)
		}()
	}
	go Flush()
	wg.Wait()
	Flush()
	if s := isAvail(); s != "unavailable" {
		t.Fatalf("expected unavailable, got %s\n", s)
	}
}

// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -run='^TestCleanupRoutine1$'
func TestCleanupRoutine1(t *testing.T) {
	go CleanupRoutine(4*time.Second, 3)
	defer ShutdownCleanupRoutine()
	t.Run("TestSet3", TestSet3)
	isAvail := func() string {
		if Check(bookseat, dfTicket) {
			return "avail"
		} else {
			return "unavailable"
		}
	}
	if s := isAvail(); s != "avail" {
		t.Fatalf("expected avail, got %s\n", s)
	}
	time.Sleep(5 * time.Second)
	if s := isAvail(); s != "unavailable" {
		t.Fatalf("expected unavailable, got %s\n", s)
	}
}

// go test -v -benchtime=10s -failfast -cpu=4 -race -benchmem -bench='^BenchmarkSet1$' -run='notmatch'
// goos: darwin
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-6567U CPU @ 3.30GHz
// BenchmarkSet1
// BenchmarkSet1-4          1000000             11840 ns/op             304 B/op          3 allocs/op
func BenchmarkSet1(b *testing.B) {
	var (
		concurrentReq int             = b.N
		wg            *sync.WaitGroup = &sync.WaitGroup{}
	)
	wg.Add(concurrentReq)
	for i := 0; i < concurrentReq; i++ {
		go func() {
			defer wg.Done()
			Set(bookseat, &ticket{film: "batman", room: 2, seat: 7})
		}()
	}
	wg.Wait()
}
