package mcache

import (
	"sync"
	"testing"
	"time"
)

// These tests must be run separately (different process).
// If not, it will use the same memory, and interfere each test case.
// If you want to run all these test cases, run test(.)sh.

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
	println(data)    // 0xc000070300
	println(data[k]) // 0x0 == nil
	data[k] = new(cacheVal)
	println(data[k].mu) // 0x0 == nil
	data[k].mu = &sync.RWMutex{}
	println(data[k].v) //  0x0 == nil
	data[k].v = make(map[any]int64)
	v := ticket{}
	println(data[k].v[v]) // 0
	data[k].v[v] = time.Now().Unix()
}

// [Set] with value as pointer.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -race -run='^TestSet1$'
func TestSet1(t *testing.T) {
	var (
		concurrent int32           = 100
		wg         *sync.WaitGroup = &sync.WaitGroup{}
		mu         *sync.Mutex     = &sync.Mutex{}
		statuses   []bool          = make([]bool, concurrent)
		maxCache   uint32          = 1 // Should only 1 cache that has been succeded to Set (returns true).
	)
	wg.Add(int(concurrent))
	for i := int32(0); i < concurrent; i++ {
		go func() {
			mu.Lock()
			defer wg.Done()
			defer mu.Unlock()
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
		concurrent int32           = 100
		wg         *sync.WaitGroup = &sync.WaitGroup{}
		mu         *sync.Mutex     = &sync.Mutex{}
		statuses   []bool          = make([]bool, concurrent)
		maxCache   uint32          = 1 // Should only 1 cache that has been succeded to Set (returns true).
	)
	wg.Add(int(concurrent))
	for i := int32(0); i < concurrent; i++ {
		go func() {
			mu.Lock()
			defer wg.Done()
			defer mu.Unlock()
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
		concurrent int32           = 100
		wg         *sync.WaitGroup = &sync.WaitGroup{}
		mu         *sync.Mutex     = &sync.Mutex{}
		statuses   []bool          = make([]bool, concurrent)
		maxCache   uint32          = 1 // Should only 1 cache that has been succeded to Set (returns true).
	)
	wg.Add(int(concurrent + 1))
	go func() {
		mu.Lock()
		defer wg.Done()
		defer mu.Unlock()
		statuses = append(statuses, Set(bookseat, &ticket{film: "batman", room: 2, seat: 7}))
	}()
	for i := int32(0); i < concurrent; i++ {
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
	var (
		concurrent int             = 1000
		mu         *sync.Mutex     = &sync.Mutex{}
		wg         *sync.WaitGroup = &sync.WaitGroup{}
	)
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func() {
			mu.Lock()
			defer wg.Done()
			defer mu.Unlock()
			Flush()
		}()
	}
	wg.Wait()
}

// Simple scenario for flush test case.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -run='^TestFlush2$'
func TestFlush2(t *testing.T) {
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

// [Set] and [Flush] run concurrently.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -run='^TestFlush3$'
func TestFlush3(t *testing.T) {
	var (
		concurrent int32           = 1000
		wg         *sync.WaitGroup = &sync.WaitGroup{}
	)
	isAvail := func() string {
		if Check(bookseat, dfTicket) {
			return "avail"
		} else {
			return "unavailable"
		}
	}
	wg.Add(int(concurrent))
	for i := int32(0); i < concurrent; i++ {
		go func() {
			defer wg.Done()
			Set(bookseat, dfTicket)
		}()
	}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			Flush()
		}()
	}
	wg.Wait()
	Flush()
	if s := isAvail(); s != "unavailable" {
		t.Fatalf("expected unavailable, got %s\n", s)
	}
}

// [Delete] run concurrently.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -run='^TestDelete1$'
func TestDelete1(t *testing.T) {
	var (
		concurrent int32           = 1000
		wg         *sync.WaitGroup = &sync.WaitGroup{}
	)
	isAvail := func() string {
		if Check(bookseat, dfTicket) {
			return "avail"
		} else {
			return "unavailable"
		}
	}
	wg.Add(int(concurrent))
	for i := 0; i < int(concurrent); i++ {
		go func() {
			defer wg.Done()
			Delete(bookseat)
		}()
	}
	wg.Wait()
	Delete(bookseat)
	if s := isAvail(); s != "unavailable" {
		t.Fatalf("expected unavailable, got %s\n", s)
	}
}

// [Set] and [Delete] run concurrently.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -run='^TestDelete2$'
func TestDelete2(t *testing.T) {
	var (
		concurrent int32           = 1000
		wg         *sync.WaitGroup = &sync.WaitGroup{}
	)
	isAvail := func() string {
		if Check(bookseat, dfTicket) {
			return "avail"
		} else {
			return "unavailable"
		}
	}
	wg.Add(int(concurrent))
	for i := int32(0); i < concurrent; i++ {
		go func() {
			defer wg.Done()
			Set(bookseat, dfTicket)
		}()
	}
	wg.Add(int(concurrent))
	for i := 0; i < int(concurrent); i++ {
		go func() {
			defer wg.Done()
			Delete(bookseat)
		}()
	}
	wg.Wait()
	Delete(bookseat)
	if s := isAvail(); s != "unavailable" {
		t.Fatalf("expected unavailable, got %s\n", s)
	}
}

// [Check] and [Delete] run concurrently.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -run='^TestDelete3$'
func TestDelete3(t *testing.T) {
	var (
		concurrent int32           = 1000
		wg         *sync.WaitGroup = &sync.WaitGroup{}
	)
	wg.Add(int(concurrent))
	for i := int32(0); i < concurrent; i++ {
		go func() {
			defer wg.Done()
			Set(bookseat, dfTicket)
		}()
	}
	wg.Add(int(concurrent))
	for i := 0; i < int(concurrent); i++ {
		go func() {
			defer wg.Done()
			Check(bookseat, dfTicket)
		}()
	}
	wg.Wait()
}

// [Set] and [Check] run concurrently.
// -count=n, n should be 1, if bigger it will fail.
// This happens because test function being runned again, but the memory of [data] remain same.
// go test -v -count=1 -failfast -cpu=4 -run='^TestCheck1$'
func TestCheck1(t *testing.T) {
	var (
		concurrent int32           = 1000
		wg         *sync.WaitGroup = &sync.WaitGroup{}
	)
	isAvail := func() string {
		if Check(bookseat, dfTicket) {
			return "avail"
		} else {
			return "unavailable"
		}
	}
	wg.Add(int(concurrent))
	for i := int32(0); i < concurrent; i++ {
		go func() {
			defer wg.Done()
			Set(bookseat, dfTicket)
		}()
	}
	wg.Add(int(concurrent))
	for i := 0; i < int(concurrent); i++ {
		go func() {
			defer wg.Done()
			Check(bookseat, dfTicket)
			if s := isAvail(); s != "avail" {
				panic("expected avail, got unavailable")
			}
		}()
	}
	wg.Wait()
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
	time.Sleep(4500 * time.Millisecond)
	if s := isAvail(); s != "unavailable" {
		t.Fatalf("expected unavailable, got %s\n", s)
	}
}

// go test -v -benchtime=10s -failfast -cpu=4 -race -benchmem -bench='^BenchmarkSet1$' -run='notmatch'
// goos: darwin
// goarch: amd64
// cpu: Intel(R) Core(TM) i7-6567U CPU @ 3.30GHz
// BenchmarkSet1
// BenchmarkSet1-4          2672319              7976 ns/op              87 B/op          3 allocs/op
func BenchmarkSet1(b *testing.B) {
	var (
		concurrent int             = b.N
		wg         *sync.WaitGroup = &sync.WaitGroup{}
	)
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func() {
			defer wg.Done()
			Set(bookseat, &ticket{film: "batman", room: 2, seat: 7})
		}()
	}
	wg.Wait()
}
