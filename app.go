package reqlock

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var (
	ticketSvc *TicketService = &TicketService{}
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	var (
		scBook     bool
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
	scBook = ticketSvc.Book(t)
	w.Header().Set("Content-Type", "application/json")
	if !scBook {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{
			"meta": map[string]any{
				"code":    400,
				"message": "Seat has been booked by someone else.",
			},
			"data": t,
		})
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"meta": map[string]any{
			"code":    200,
			"message": "Ticket secured",
		},
		"data": t,
	})
}
