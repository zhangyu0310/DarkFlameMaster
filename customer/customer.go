package customer

import (
	"DarkFlameMaster/config"
	"DarkFlameMaster/ticket"
	"errors"
	zlog "github.com/zhangyu0310/zlogger"
	"time"
)

var (
	customerInfo map[string]*Customer
	strategy     Strategy
)

var (
	ErrCustomerNotExist = errors.New("customer not exist")
)

func Init() error {
	cfg := config.GetGlobalConfig()
	var reader Reader
	switch cfg.CustomerFileType {
	case config.AliPay:
		reader = &AliPayCustomerInfoReader{}
	case config.WeChat:
	case config.TestCF:
	default:
		reader = &TestCustomerInfoReader{}
	}
	err := reader.Init(cfg.CustomerFile)
	if err != nil {
		zlog.Error("Init customer file reader failed, err:", err)
		return err
	}
	customers, err := reader.Read()
	if err != nil {
		zlog.Error("Read customer info failed, err:", err)
		return err
	}
	customerInfo = customers

	switch cfg.ChooseSeatStrategy {
	case config.PayTimeOneByOne:
		strategy = &PayTimeOneByOne{}
	case config.NoLimit:
		strategy = &NoLimit{}
	case config.TestCS:
	default:
		strategy = &NoLimit{}
	}
	err = strategy.Init(customerInfo)
	if err != nil {
		zlog.Error("Init choose seat strategy failed, err:", err)
		return err
	}
	return nil
}

func GetCustomer(proof string) (*Customer, error) {
	cus, ok := customerInfo[proof]
	if !ok {
		return nil, ErrCustomerNotExist
	}
	return cus, nil
}

func CanChooseSeat(cus *Customer) bool {
	return strategy.CanChooseSeat(cus)
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
