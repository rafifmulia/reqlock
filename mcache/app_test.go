package mcache

import (
	"sync"
	"testing"
	"time"
)

type Ticket struct {
	Id, Room, Seat int32
	Film           string
}

var (
	ck       string = "bookseat" // cacheKey
	dfTicket Ticket = Ticket{Film: "batman", Room: 2, Seat: 7}
)

// go test -v -count=1 -failfast -cpu=1 -run='^TestDefaultValue$'
func TestDefaultValue(t *testing.T) {
	k := ck
	println(data)    // 0xc0000a2300
	println(data[k]) // 0x0 == nil
	data[k] = new(objKey)
	println(data[k].mtm) // 0
	data[k].mtm = time.Now().Unix()
	println(data[k].v) // 0x0 == nil
	v := Ticket{}
	data[k].v = make(map[any]bool)
	data[k].v[v] = true
	println(data[k])
}

// go test -v -count=1 -failfast -cpu=4 -race -run='^TestSet1$'
func TestSet1(t *testing.T) {
	var (
		concurrentReq int32           = 40
		wg            *sync.WaitGroup = &sync.WaitGroup{}
	)
	ticket := &Ticket{Film: "batman", Room: 2, Seat: 7}
	Set(ck, ticket)
	wg.Add(int(concurrentReq))
	for i := int32(0); i < concurrentReq; i++ {
		go func() {
			defer wg.Done()
			ticket := &Ticket{Film: "batman", Room: 2, Seat: 7}
			status := Set(ck, ticket)
			if status {
				panic("duplicated bookseat")
			}
		}()
	}
	wg.Wait()
}

// go test -v -count=1 -failfast -cpu=4 -race -run='^TestSet2$'
func TestSet2(t *testing.T) {
	var (
		concurrentReq int32           = 40
		wg            *sync.WaitGroup = &sync.WaitGroup{}
	)
	Set(ck, dfTicket)
	wg.Add(int(concurrentReq))
	for i := int32(0); i < concurrentReq; i++ {
		go func() {
			defer wg.Done()
			ticket := Ticket{Film: "batman", Room: 2, Seat: 7}
			status := Set(ck, ticket)
			if status {
				panic("duplicated bookseat")
			}
		}()
	}
	wg.Wait()
}

// go test -v -count=1 -failfast -cpu=4 -race -run='^TestSet3$'
func TestSet3(t *testing.T) {
	var (
		concurrentReq int32           = 40
		wg            *sync.WaitGroup = &sync.WaitGroup{}
	)
	Set(ck, dfTicket)
	wg.Add(int(concurrentReq))
	for i := int32(0); i < concurrentReq; i++ {
		go func() {
			defer wg.Done()
			ticket := &Ticket{Film: "batman", Room: 2, Seat: 7}
			status := Set(ck, ticket)
			if status {
				panic("duplicated bookseat")
			}
		}()
	}
	wg.Wait()
}

// go test -v -count=5 -failfast -cpu=4 -run='^TestCleanupRoutine1$'
func TestCleanupRoutine1(t *testing.T) {
	go CleanupRoutine(4*time.Second, 3)
	t.Run("TestSet3", TestSet3)
	isAvail := func() string {
		if Check(ck, dfTicket) {
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
// BenchmarkSet1-4          2622802              6384 ns/op              89 B/op          3 allocs/op
func BenchmarkSet1(b *testing.B) {
	var (
		concurrentReq int             = b.N
		wg            *sync.WaitGroup = &sync.WaitGroup{}
	)
	ticket := &Ticket{Film: "batman", Room: 2, Seat: 7}
	Set(ck, ticket)
	wg.Add(concurrentReq)
	for i := 0; i < concurrentReq; i++ {
		go func() {
			defer wg.Done()
			ticket := &Ticket{Film: "batman", Room: 2, Seat: 7}
			status := Set(ck, ticket)
			if status {
				panic("duplicated bookseat")
			}
		}()
	}
	wg.Wait()
}
