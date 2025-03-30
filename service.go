package reqlock

import (
	"math/rand"
	"sync"
	"time"
)

type TicketService struct {
	tickets []*Ticket
	mux     *sync.Mutex
}

func (svc *TicketService) Book(t *Ticket) bool {
	svc.mux.Lock()
	defer svc.mux.Unlock()
	time.Sleep(500 * time.Millisecond) // Assume query to database takes 500 milisecond.
	for i := 0; i < len(svc.tickets); i++ {
		if svc.tickets[i].Film == t.Film && svc.tickets[i].Room == t.Room && svc.tickets[i].Seat == t.Seat {
			return false
		}
	}
	svc.tickets = append(svc.tickets, &Ticket{Id: rand.Int31(), Film: t.Film, Room: t.Room, Seat: t.Seat})
	return true
}
