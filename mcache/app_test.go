package mcache

import (
	"sync"
	"testing"
)

type Ticket struct {
	Id, Room, Seat int32
	Film           string
}

// go test -v -count=1 -failfast -cpu=4 -run='^TestSet1$'
func TestSet1(t *testing.T) {
	var (
		concurrentReq int32           = 40
		wg            *sync.WaitGroup = &sync.WaitGroup{}
	)
	ticket := &Ticket{Film: "batman", Room: 2, Seat: 7}
	Set("bookseat", ticket)
	wg.Add(int(concurrentReq))
	for i := int32(0); i < concurrentReq; i++ {
		go func() {
			defer wg.Done()
			status := Set("bookseat", ticket)
			if status {
				panic("duplicated bookseat")
			}
		}()
	}
	wg.Wait()
}
