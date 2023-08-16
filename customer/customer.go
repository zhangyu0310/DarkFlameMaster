package customer

import (
	"DarkFlameMaster/config"
	"DarkFlameMaster/ticket"
	"errors"
	zlog "github.com/zhangyu0310/zlogger"
	"time"
)

var (
	ErrCustomerNotExist  = errors.New("customer not exist")
	ErrStrategicConflict = errors.New("customer strategic conflict")
)

var cusMgr *Manager

type Manager struct {
	reader   Reader
	strategy Strategy
}

func (mgr *Manager) Init() error {
	cfg := config.GetGlobalConfig()
	var reader Reader
	switch cfg.CustomerType {
	case config.AliPay:
		reader = &AliPayCustomerInfoReader{}
	case config.WeChat:
	case config.NoPay:
		reader = &NoPayCustomerInfoReader{}
	case config.TestCF:
		reader = &TestCustomerInfoReader{}
	default:
		reader = &TestCustomerInfoReader{}
	}
	err := reader.Init(cfg.CustomerFile)
	if err != nil {
		zlog.Error("Init customer file reader failed, err:", err)
		return err
	}
	cus, err := reader.Read()
	if err != nil {
		zlog.Error("Read customer info failed, err:", err)
		return err
	}
	mgr.reader = reader

	switch cfg.ChooseSeatStrategy {
	case config.PayTimeOneByOne:
		mgr.strategy = &PayTimeOneByOne{}
	case config.NoLimit:
		mgr.strategy = &NoLimit{}
	default:
		mgr.strategy = &NoLimit{}
	}
	err = mgr.strategy.Init(cus)
	if err != nil {
		zlog.Error("Init choose seat strategy failed, err:", err)
		return err
	}
	return nil
}

func (mgr *Manager) GetCustomer(proof string) (*Customer, error) {
	return mgr.reader.GetCustomerInfo(proof)
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

func (c *Customer) RemainTicketNumber() uint {
	return c.TicketNum - uint(len(c.Tickets))
}
