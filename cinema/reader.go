package cinema

import (
	"DarkFlameMaster/seat"
	"math/rand"
	"sync"
)

type Reader interface {
	Init(string) error
	Read() ([][]*seat.Seat, uint, uint, error)
}

type JsonTypeSeatInfoReader struct {
	infoFilePath string
}

func (r *JsonTypeSeatInfoReader) Init(seatFilePath string) error {
	r.infoFilePath = seatFilePath
	return nil
}

// TODO: 实现Json格式的座位表读取
func (r *JsonTypeSeatInfoReader) Read() ([][]*seat.Seat, uint, uint, error) {
	return nil, 0, 0, nil
}

type TestSeatInfoReader struct{}

func (r *TestSeatInfoReader) Init(string) error {
	return nil
}

func (r *TestSeatInfoReader) Read() ([][]*seat.Seat, uint, uint, error) {
	seats := make([][]*seat.Seat, 0)
	for i := uint(0); i < 10; i++ {
		seats = append(seats, make([]*seat.Seat, 0))
		for j := uint(0); j < 10; j++ {
			seats[i] = append(seats[i], &seat.Seat{
				Mutex:  sync.Mutex{},
				Row:    i + 1,
				Column: j + 1,
				Status: seat.Available,
			})
		}
	}
	for i := 0; i < 20; i++ {
		r := rand.Intn(9)
		c := rand.Intn(9)
		seats[r][c].Status = seat.Elected
	}
	return seats, 10, 10, nil
}
