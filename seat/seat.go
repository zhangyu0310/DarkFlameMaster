package seat

import (
	"DarkFlameMaster/ticket"
	"sync"
)

type Status int

const (
	Hit Status = iota
	Available
	Elected
	Invalid
)

func Status2Str(status Status) string {
	switch status {
	case Available:
		return "Available"
	case Elected:
		return "Elected"
	case Invalid:
		return "Invalid"
	default:
		return "Invalid"
	}
}

type Seat struct {
	sync.Mutex

	Row    uint   `json:"row"`
	Column uint   `json:"col"`
	Status Status `json:"status"`

	Ticket *ticket.Ticket
}

func (s *Seat) IsAvailable() bool {
	return s.Status == Available
}

func (s *Seat) BandageTicket(t *ticket.Ticket) {
	s.Status = Elected
	s.Ticket = t
}
