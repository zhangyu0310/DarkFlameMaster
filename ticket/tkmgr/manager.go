package tkmgr

import (
	"DarkFlameMaster/cinema"
	"DarkFlameMaster/config"
	"DarkFlameMaster/customer"
	"DarkFlameMaster/seat"
	"DarkFlameMaster/storage"
	"DarkFlameMaster/ticket"
	"errors"
	zlog "github.com/zhangyu0310/zlogger"
)

var (
	ErrSeatIsNotAvailable = errors.New("seat is not available")
)

var mgr *TicketManager

type TicketManager struct {
	tickets []*ticket.Ticket

	Storage storage.TicketStorage
}

func Init() error {
	mgr = &TicketManager{}
	return mgr.init()
}

func MakeTickets(c *customer.Customer, seats []*seat.Seat) ([]*ticket.Ticket, error) {
	return mgr.makeTickets(c, seats)
}

func (m *TicketManager) init() error {
	cfg := config.GetGlobalConfig()
	switch cfg.DbType {
	case config.LevelDB:
		m.Storage = &storage.TicketLevelDB{}
	case config.MySQL:
	case config.ZDB:
	default:
		m.Storage = &storage.TicketLevelDB{}
	}
	err := m.Storage.Init(cfg.DbPath)
	if err != nil {
		zlog.Error("Init ticket storage failed, err:", err)
		return err
	}

	// Read all tickets info from database, then connect them with other modules.
	m.tickets, err = m.Storage.ReadAll()
	if err != nil {
		zlog.Error("Read all tickets from db failed, err:", err)
		return err
	}
	for _, t := range m.tickets {
		zlog.DebugF("Read ticket [%s] from db. [%v]", t.ID, t)
		// Bandage ticket to seat
		s, err := cinema.GetSeat(t.Row, t.Column)
		if err != nil {
			zlog.Error("Get seat failed, err:", err)
			return err
		}
		s.BandageTicket(t)
		// Bandage ticket to customer
		c, err := customer.GetCustomer(t.CustomerProof)
		if err != nil {
			zlog.Error("Get customer failed, err:", err)
			return err
		}
		c.AddTicket(t)
	}
	return nil
}

func unlockAllSeats(seats []*seat.Seat, target *seat.Seat) {
	for _, s := range seats {
		s.Unlock()
		if s == target {
			break
		}
	}
}

func (m *TicketManager) makeTickets(customer *customer.Customer,
	seats []*seat.Seat) ([]*ticket.Ticket, error) {
	// Lock seats and check its status.
	for _, s := range seats {
		s.Lock()
		if !s.IsAvailable() {
			zlog.DebugF("Seat [%v] is not available", s)
			unlockAllSeats(seats, s)
			return nil, ErrSeatIsNotAvailable
		}
	}
	defer unlockAllSeats(seats, nil)
	// All seats are available and locked, make tickets.
	tickets := make([]*ticket.Ticket, 0, len(seats))
	for _, s := range seats {
		t := ticket.NewTicket(customer.Proof, s.Row, s.Column)
		tickets = append(tickets, t)
	}
	// Save tickets to database.
	err := m.Storage.WriteTickets(tickets)
	if err != nil {
		zlog.Error("Write tickets to db failed, err:", err)
		return nil, err
	}
	// Bandage tickets to seats & save them to memory.
	for i, t := range tickets {
		seats[i].BandageTicket(t)
		m.tickets = append(m.tickets, t)
	}
	return tickets, nil
}

func (m *TicketManager) CheckTickets(proof string) []*ticket.Ticket {
	t := make([]*ticket.Ticket, 0)
	for _, v := range m.tickets {
		if v.CustomerProof == proof {
			t = append(t, v)
		}
	}
	return t
}

func Dump() ([]*ticket.Ticket, error) {
	return mgr.Storage.ReadAll()
}

func CheckTickets(proof string) []*ticket.Ticket {
	return mgr.CheckTickets(proof)
}
