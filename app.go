package reqlock

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

var (
	ticketSvc *TicketService = &TicketService{mux: &sync.Mutex{}}
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	var (
		booked     bool
		room, seat int
		film       string
		err        error
	)
	err = r.ParseForm()
	if err != nil {
		panic(err)
	}
	film = r.PostFormValue("film")
	room, err = strconv.Atoi(r.PostFormValue("room"))
	if err != nil {
		panic(err)
	}
	seat, err = strconv.Atoi(r.PostFormValue("seat"))
	if err != nil {
		panic(err)
	}
	t := &Ticket{Film: film, Room: int32(room), Seat: int32(seat)}
	booked = ticketSvc.Book(t)
	w.Header().Set("Content-Type", "application/json")
	if !booked {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{
			"meta": map[string]any{
				"code":    400,
				"message": "Seat has been booked by someone else.",
			},
		})
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"meta": map[string]any{
			"code":    200,
			"message": "Ticket secured",
		},
	})
}
