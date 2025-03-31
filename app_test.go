package reqlock

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rafifmulia/reqlock/mcache"
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
	go mcache.CleanupRoutine(4*time.Second, 2)
	defer mcache.ShutdownCleanupRoutine()
	wg.Add(int(concurrentReq))
	for i := int32(0); i < concurrentReq; i++ {
		go doReq()
	}
	wg.Wait()
	t.Log("Total Concurrent Request:", concurrentReq, "Total Booked:", len(ticketSvc.tickets))
	if len(ticketSvc.tickets) > 1 {
		for _, v := range ticketSvc.tickets {
			t.Error(v.Id, v.Film, v.Room, v.Seat)
		}
	}
}

// go test -v -count=1 -failfast -cpu=4 -run='^TestRequestHandler2$'
func TestRequestHandler2(t *testing.T) {
	var (
		concurrentReq int32           = 1000
		wg            *sync.WaitGroup = &sync.WaitGroup{}
	)
	ticketSvc.tickets = make([]*Ticket, 0, 1)
	doReq := func() {
		defer wg.Done()
		room := rand.Int31n(10)
		seat := rand.Int31n(10)
		py := strings.NewReader(fmt.Sprintf("film=batman&room=%d&seat=%d", room, seat))
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
	t.Log("Total Concurrent Request:", concurrentReq, "Total Booked:", len(ticketSvc.tickets))
	if len(ticketSvc.tickets) > 1 {
		type ticketKey struct {
			film string
			room int32
			seat int32
		}
		seen := make(map[ticketKey]bool)
		for _, v := range ticketSvc.tickets {
			key := ticketKey{
				film: v.Film,
				room: v.Room,
				seat: v.Seat,
			}
			if seen[key] {
				t.Fatalf("Duplicate ticket found: Film=%s, Room=%d, Seat=%d", v.Film, v.Room, v.Seat)
			}
			seen[key] = true
		}
	}
}

// This fuzz test run as a worker, it means run as different process.
// If found duplicated ticket, this is because every process has its own memory.
// No need to wory about this, because in TestRequestHandler* there aren't found duplicated tickets.
// go test -v -fuzztime=1m -cpu=1 -fuzz='^FuzzRequestHandler1$' -run='notmatch'
func FuzzRequestHandler1(f *testing.F) {
	var (
		of  *os.File
		err error
	)
	of, err = os.OpenFile("testdata/fuzz.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0o660)
	defer of.Close()
	if err != nil {
		panic(err)
	}
	ticketSvc.tickets = make([]*Ticket, 0, 1)
	tc := []uint8{1, 2, 3}
	for _, v := range tc {
		f.Add(v, v)
	}
	f.Fuzz(func(t *testing.T, room, seat uint8) {
		py := strings.NewReader(fmt.Sprintf("film=batman&room=%d&seat=%d", room, seat))
		req := httptest.NewRequest("POST", "/ticket/book", py)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		RequestHandler(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()
		a := []byte(t.Name())
		b, _ := io.ReadAll(resp.Body)
		c := append(a, b...)
		of.Write(c)
	})
}
