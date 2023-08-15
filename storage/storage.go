package storage

import (
	"DarkFlameMaster/ticket"
)

type TicketStorage interface {
	Init(string, ...interface{}) error
	ReadAll() ([]*ticket.Ticket, error)
	WriteTickets([]*ticket.Ticket) error
}
