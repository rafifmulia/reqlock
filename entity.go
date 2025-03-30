package reqlock

type Ticket struct {
	Id, Room, Seat int32
	Film           string
}
