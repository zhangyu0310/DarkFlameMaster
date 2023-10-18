package customer

import (
	"DarkFlameMaster/config"
	"DarkFlameMaster/ticket"
	"errors"
	"time"

	zlog "github.com/zhangyu0310/zlogger"
)

var (
	ErrCustomerNotExist  = errors.New("customer not exist")
	ErrStrategicConflict = errors.New("customer strategic conflict")
)

var cusMgr *Manager

type Manager struct {
	storage  Storage
	strategy Strategy
}

func (mgr *Manager) Init() error {
	cfg := config.GetGlobalConfig()
	var storage Storage
	switch cfg.CustomerType {
	case config.AliPay:
		storage = &AliPayCustomerInfoStorage{}
	case config.WeChat:
		panic("implement me")
	case config.QQNum:
		storage = &QQNumCustomerInfoStorage{}
	case config.NoPay:
		storage = &NoPayCustomerInfoStorage{}
	case config.TestCF:
		storage = &TestCustomerInfoStorage{}
	default:
		storage = &TestCustomerInfoStorage{}
	}
	err := storage.Init(cfg.CustomerFile)
	if err != nil {
		zlog.Error("Init customer file reader failed, err:", err)
		return err
	}
	customer, err := storage.Read()
	if err != nil {
		zlog.Error("Read customer info failed, err:", err)
		return err
	}
	mgr.storage = storage
	if cfg.RootUserName != "" {
		_ = mgr.storage.AddCustomer(&Customer{
			Proof:     cfg.RootUserName,
			PayTime:   time.Now(),
			TicketNum: 2333,
			Tickets:   nil,
		})
	}

	switch cfg.ChooseSeatStrategy {
	case config.PayTimeOneByOne:
		mgr.strategy = &PayTimeOneByOne{}
	case config.NoLimit:
		mgr.strategy = &NoLimit{}
	default:
		mgr.strategy = &NoLimit{}
	}
	err = mgr.strategy.Init(customer)
	if err != nil {
		zlog.Error("Init choose seat strategy failed, err:", err)
		return err
	}
	return nil
}

func (mgr *Manager) GetCustomer(proof string) (*Customer, error) {
	return mgr.storage.GetCustomerInfo(proof)
}

func (mgr *Manager) CanChooseSeat(cus *Customer) bool {
	return mgr.strategy.CanChooseSeat(cus)
}

func Init() error {
	cusMgr = &Manager{}
	return cusMgr.Init()
}

func GetCustomer(proof string) (*Customer, error) {
	return cusMgr.GetCustomer(proof)
}

func CanChooseSeat(cus *Customer) bool {
	return cusMgr.CanChooseSeat(cus)
}

type Customer struct {
	Proof     string
	PayTime   time.Time
	TicketNum uint
	Tickets   []*ticket.Ticket
}

func (c *Customer) AddTicket(t *ticket.Ticket) {
	c.Tickets = append(c.Tickets, t)
}

func (c *Customer) DeleteTicket(t *ticket.Ticket) {
	for i, v := range c.Tickets {
		if v.ID == t.ID {
			c.Tickets = append(c.Tickets[:i], c.Tickets[i+1:]...)
			return
		}
	}
}

func (c *Customer) RemainTicketNumber() uint {
	return c.TicketNum - uint(len(c.Tickets))
}
