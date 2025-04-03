package reqlock

import (
	"math/rand"
	"sync"
	"time"

	"github.com/rafifmulia/reqlock/mcache"
)

const (
	keyCache = "bookseat"
)

type TicketService struct {
	tickets []*Ticket
	mu      *sync.Mutex
}

func (svc *TicketService) Book(t *Ticket) bool {
	if !mcache.Set(keyCache, t) {
		return false
	}
	time.Sleep(200 * time.Millisecond) // Assume insert data to database takes 200 miliseconds.
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.tickets = append(svc.tickets, &Ticket{Id: rand.Int31(), Film: t.Film, Room: t.Room, Seat: t.Seat})
	return true
}
