package customer

import (
	"fmt"
	"github.com/google/uuid"
	zlog "github.com/zhangyu0310/zlogger"
	"os"
	"time"
)

type Reader interface {
	Init(string) error
	Read() (map[string]*Customer, error)
	GetCustomerInfo(string) (*Customer, error)
}

type AliPayCustomerInfoReader struct {
	infoFilePath string
	customers    map[string]*Customer
}

func (r *AliPayCustomerInfoReader) Init(customerFilePath string) error {
	r.infoFilePath = customerFilePath
	r.customers = make(map[string]*Customer)
	return nil
}

// TODO: 实现读取支付宝账单信息
func (r *AliPayCustomerInfoReader) Read() (map[string]*Customer, error) {
	data, err := os.ReadFile(r.infoFilePath)
	if err != nil {
		zlog.Error("Read customer info file failed, err:", err)
		return nil, err
	}
	fmt.Println(string(data))
	return nil, nil
}

func (r *AliPayCustomerInfoReader) GetCustomerInfo(proof string) (*Customer, error) {
	cus, ok := r.customers[proof]
	if !ok {
		return nil, ErrCustomerNotExist
	}
	return cus, nil
}

type TestCustomerInfoReader struct {
	customers map[string]*Customer
}

func (r *TestCustomerInfoReader) Init(string) error {
	r.customers = make(map[string]*Customer)
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

func (r *TestCustomerInfoReader) GetCustomerInfo(proof string) (*Customer, error) {
	cus, ok := r.customers[proof]
	if !ok {
		return nil, ErrCustomerNotExist
	}
	return cus, nil
}

type NoPayCustomerInfoReader struct {
	customers map[string]*Customer
}

func (r *NoPayCustomerInfoReader) Init(string) error {
	r.customers = make(map[string]*Customer)
	return nil
}

func (r *NoPayCustomerInfoReader) Read() (map[string]*Customer, error) {
	return nil, nil
}

func (r *NoPayCustomerInfoReader) GetCustomerInfo(proof string) (*Customer, error) {
	cus, ok := r.customers[proof]
	if !ok {
		cus = &Customer{
			Proof:     proof,
			PayTime:   time.Now(),
			TicketNum: 2333,
		}
		r.customers[proof] = cus
	}
	return cus, nil
}
