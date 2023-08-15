package customer

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Reader interface {
	Init(string) error
	Read() (map[string]*Customer, error)
}

type AliPayCustomerInfoReader struct {
	infoFilePath string
}

func (r *AliPayCustomerInfoReader) Init(customerFilePath string) error {
	r.infoFilePath = customerFilePath
	return nil
}

// TODO: 实现读取支付宝账单信息
func (r *AliPayCustomerInfoReader) Read() (map[string]*Customer, error) {
	return nil, nil
}

type TestCustomerInfoReader struct{}

func (r *TestCustomerInfoReader) Init(string) error {
	return nil
}

func (r *TestCustomerInfoReader) Read() (map[string]*Customer, error) {
	customers := make(map[string]*Customer)
	for i := 0; i < 50; i++ {
		id := uuid.NewString()
		fmt.Println(id)
		customers[id] = &Customer{
			Proof:     id,
			PayTime:   time.Now(),
			TicketNum: 1,
			Tickets:   nil,
		}
	}

	return customers, nil
}
