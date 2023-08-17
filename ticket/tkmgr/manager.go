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
	"sync"
)

var (
	ErrSeatIsNotAvailable = errors.New("seat is not available")
)

var mgr *TicketManager

type TicketManager struct {
	sync.Mutex
	tickets []*ticket.Ticket

	Storage storage.TicketStorage
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
		customer.AddTicket(t)
	}
	m.Lock()
	m.tickets = append(m.tickets, tickets...)
	m.Unlock()
	return tickets, nil
}

func (m *TicketManager) CheckTickets(proof string) []*ticket.Ticket {
	cus, _ := customer.GetCustomer(proof)
	return cus.Tickets
}

func (m *TicketManager) DeleteTicketsByIds(ids []string) error {
	m.Lock()
	// 将当前所有的票据信息保存到map中，方便后续查找
	tkMap := make(map[string]*ticket.Ticket)
	for _, v := range m.tickets {
		tkMap[v.ID] = v
	}
	// 收集需要删除的票据信息
	delTk := make([]*ticket.Ticket, 0, len(ids))
	for _, id := range ids {
		tk, ok := tkMap[id]
		if ok {
			delTk = append(delTk, tk)
			delete(tkMap, id)
		}
	}
	// 将剩余的票据信息保存到内存
	m.tickets = make([]*ticket.Ticket, 0, len(tkMap))
	for _, v := range tkMap {
		m.tickets = append(m.tickets, v)
	}
	// 清理数据库中的数据和用户、座位关联数据
	for _, v := range delTk {
		cus, _ := customer.GetCustomer(v.CustomerProof)
		cus.DeleteTicket(v)
		s, _ := cinema.GetSeat(v.Row, v.Column)
		s.UnBandageTicket()
	}
	err := m.Storage.DeleteTickets(delTk)
	if err != nil {
		zlog.Error("Delete tickets from db failed, err:", err)
	}
	m.Unlock()
	return err
}

func (m *TicketManager) DeleteTicketsBySeats(seats []*seat.Seat) error {
	ids := make([]string, 0, len(seats))
	for _, v := range seats {
		s, err := cinema.GetSeat(v.Row, v.Column)
		if err != nil {
			zlog.Error("Get seat failed, err:", err)
			return err
		}
		if s.Ticket == nil {
			zlog.Error("Seat is not bandage to ticket")
			return errors.New("seat is not bandage to ticket")
		}
		ids = append(ids, s.Ticket.ID)
	}
	return m.DeleteTicketsByIds(ids)
}

func Init() error {
	mgr = &TicketManager{}
	return mgr.init()
}

func MakeTickets(c *customer.Customer, seats []*seat.Seat) ([]*ticket.Ticket, error) {
	return mgr.makeTickets(c, seats)
}

func Dump() ([]*ticket.Ticket, error) {
	return mgr.Storage.ReadAll()
}

func CheckTickets(proof string) []*ticket.Ticket {
	return mgr.CheckTickets(proof)
}

func DeleteTicketsByIds(ids []string) error {
	return mgr.DeleteTicketsByIds(ids)
}

func DeleteTicketsBySeats(seats []*seat.Seat) error {
	return mgr.DeleteTicketsBySeats(seats)
}
