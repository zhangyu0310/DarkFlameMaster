package customer

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	zlog "github.com/zhangyu0310/zlogger"
)

type Storage interface {
	Init(string) error
	Read() (map[string]*Customer, error)
	GetCustomerInfo(string) (*Customer, error)
	AddCustomer(*Customer) error
}

func normalGetCustomerInfo(customers map[string]*Customer, proof string) (*Customer, error) {
	cus, ok := customers[proof]
	if !ok {
		return nil, ErrCustomerNotExist
	}
	return cus, nil
}

type AliPayCustomerInfoStorage struct {
	infoFilePath string
	customers    map[string]*Customer
}

func (r *AliPayCustomerInfoStorage) Init(customerFilePath string) error {
	r.infoFilePath = customerFilePath
	r.customers = make(map[string]*Customer)
	return nil
}

// TODO: 实现读取支付宝账单信息
func (r *AliPayCustomerInfoStorage) Read() (map[string]*Customer, error) {
	data, err := os.ReadFile(r.infoFilePath)
	if err != nil {
		zlog.Error("Read customer info file failed, err:", err)
		return nil, err
	}
	fmt.Println(string(data))
	panic("implement me")
	return nil, nil
}

func (r *AliPayCustomerInfoStorage) GetCustomerInfo(proof string) (*Customer, error) {
	return normalGetCustomerInfo(r.customers, proof)
}

func (r *AliPayCustomerInfoStorage) AddCustomer(customer *Customer) error {
	r.customers[customer.Proof] = customer
	return nil
}

type TestCustomerInfoStorage struct {
	customers map[string]*Customer
}

func (r *TestCustomerInfoStorage) Init(string) error {
	r.customers = make(map[string]*Customer)
	return nil
}

func (r *TestCustomerInfoStorage) Read() (map[string]*Customer, error) {
	for i := 0; i < 50; i++ {
		id := uuid.NewString()
		fmt.Println(id)
		r.customers[id] = &Customer{
			Proof:     id,
			PayTime:   time.Now(),
			TicketNum: 1,
			Tickets:   nil,
		}
	}

	return r.customers, nil
}

func (r *TestCustomerInfoStorage) GetCustomerInfo(proof string) (*Customer, error) {
	return normalGetCustomerInfo(r.customers, proof)
}

func (r *TestCustomerInfoStorage) AddCustomer(customer *Customer) error {
	r.customers[customer.Proof] = customer
	return nil
}

type QQNumCustomerInfoStorage struct {
	infoFilePath string
	customers    map[string]*Customer
}

func (r *QQNumCustomerInfoStorage) Init(file string) error {
	r.infoFilePath = file
	r.customers = make(map[string]*Customer)
	return nil
}

func (r *QQNumCustomerInfoStorage) Read() (map[string]*Customer, error) {
	f, err := excelize.OpenFile(r.infoFilePath)
	if err != nil {
		zlog.Error("Open customer info excel failed, err:", err)
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			zlog.Error("Close customer info excel failed, err:", err)
		}
	}()
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		zlog.Error("Get all rows from customer info excel failed, err:", err)
		return nil, err
	}
	for _, row := range rows {
		r.customers[row[0]] = &Customer{
			Proof:     row[0],
			PayTime:   time.Now(),
			TicketNum: 3, // TODO: Max ticket can be setting.
			Tickets:   nil,
		}
	}
	return nil, nil
}

func (r *QQNumCustomerInfoStorage) GetCustomerInfo(proof string) (*Customer, error) {
	return normalGetCustomerInfo(r.customers, proof)
}

func (r *QQNumCustomerInfoStorage) AddCustomer(customer *Customer) error {
	r.customers[customer.Proof] = customer
	return nil
}

type NoPayCustomerInfoStorage struct {
	customers map[string]*Customer
}

func (r *NoPayCustomerInfoStorage) Init(string) error {
	r.customers = make(map[string]*Customer)
	return nil
}

func (r *NoPayCustomerInfoStorage) Read() (map[string]*Customer, error) {
	return nil, nil
}

func (r *NoPayCustomerInfoStorage) GetCustomerInfo(proof string) (*Customer, error) {
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

func (r *NoPayCustomerInfoStorage) AddCustomer(customer *Customer) error {
	r.customers[customer.Proof] = customer
	return nil
}
