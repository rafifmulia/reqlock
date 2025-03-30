package reqlock

import (
	"math/rand"
	"time"

	"github.com/rafifmulia/reqlock/mcache"
)

const (
	keyCache = "bookseat"
)

type TicketService struct {
	tickets []*Ticket
}

func (svc *TicketService) Book(t *Ticket) bool {
	if !mcache.Set(keyCache, t) {
		return false
	}
	time.Sleep(200 * time.Millisecond) // Assume insert data to database takes 200 miliseconds.
	svc.tickets = append(svc.tickets, &Ticket{Id: rand.Int31(), Film: t.Film, Room: t.Room, Seat: t.Seat})
	return true
}
