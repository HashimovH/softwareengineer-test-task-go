package domain

type Score struct {
	Category string
	Rating   int
	Date     string
	Score    int32
}

type ScoreByTicket struct {
	TicketId int32
	Category string
	Score    int32
}
