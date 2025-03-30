package reqlock

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// go test -v -count=1 -failfast -cpu=4 -run='^TestRequestHandler1$'
func TestRequestHandler1(t *testing.T) {
	var (
		concurrentReq int32           = 40
		wg            *sync.WaitGroup = &sync.WaitGroup{}
	)
	ticketSvc.tickets = make([]*Ticket, 0, 1)
	doReq := func() {
		defer wg.Done()
		py := strings.NewReader("film=batman&room=2&seat=7")
		req := httptest.NewRequest("POST", "/ticket/book", py)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		RequestHandler(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()
		data := map[string]any{}
		json.NewDecoder(resp.Body).Decode(&data)
	}
	wg.Add(int(concurrentReq))
	for i := int32(0); i < concurrentReq; i++ {
		go doReq()
	}
	wg.Wait()
	t.Log("Total Booked:", len(ticketSvc.tickets))
	if len(ticketSvc.tickets) > 1 {
		for _, v := range ticketSvc.tickets {
			t.Error(v.Id, v.Film, v.Room, v.Seat)
		}
	}
}
