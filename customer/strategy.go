package customer

import "sort"

type Strategy interface {
	Init(map[string]*Customer) error
	CanChooseSeat(*Customer) bool
}

// PayTimeOneByOne 根据支付时间顺序依次选座
type PayTimeOneByOne struct {
	payTimeOrder []*Customer
}

func (s *PayTimeOneByOne) Init(cus map[string]*Customer) error {
	for _, v := range cus {
		s.payTimeOrder = append(s.payTimeOrder, v)
	}
	sort.Slice(s.payTimeOrder, func(i, j int) bool {
		return s.payTimeOrder[i].PayTime.Before(s.payTimeOrder[j].PayTime)
	})
	return nil
}

func (s *PayTimeOneByOne) CanChooseSeat(cus *Customer) bool {
	for _, v := range s.payTimeOrder {
		if v.RemainTicketNumber() > 0 {
			return v == cus
		}
	}
	return false
}

type NoLimit struct{}

func (s *NoLimit) Init(map[string]*Customer) error {
	return nil
}

func (s *NoLimit) CanChooseSeat(*Customer) bool {
	return true
}
