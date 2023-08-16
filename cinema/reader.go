package cinema

import (
	"DarkFlameMaster/seat"
	"encoding/json"
	"errors"
	zlog "github.com/zhangyu0310/zlogger"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Reader interface {
	Init(string) error
	Read() ([][]*seat.Seat, []*seat.Block, uint, uint, error)
}

type JsonTypeSeatInfoReader struct {
	infoFilePath string
}

type BlockInfo struct {
	Row       string `json:"row"`
	Col       uint   `json:"col"`
	Direction string `json:"direction"`
	BlockNum  uint   `json:"blockNum"`
}

type SeatInfo struct {
	Row     string `json:"row"`
	SeatNum uint   `json:"seatNum"`
}

type JsonSeat struct {
	MaxRow    uint         `json:"maxRow"`
	MaxCol    uint         `json:"maxCol"`
	BlockInfo []*BlockInfo `json:"blockInfo"`
	SeatInfo  []*SeatInfo  `json:"seatInfo"`
}

func (r *JsonTypeSeatInfoReader) Init(seatFilePath string) error {
	r.infoFilePath = seatFilePath
	return nil
}

func getStartAndEnd(str string) (uint, uint, error) {
	sp := strings.Split(str, "-")
	if len(sp) != 2 {
		return 0, 0, errors.New("invalid input, without start & end")
	}
	start, err := strconv.Atoi(sp[0])
	if err != nil {
		return 0, 0, err
	}
	end, err := strconv.Atoi(sp[1])
	if err != nil {
		return 0, 0, err
	}
	if end < start {
		return 0, 0, errors.New("end is less than start")
	}
	return uint(start), uint(end), nil
}

func (r *JsonTypeSeatInfoReader) Read() ([][]*seat.Seat, []*seat.Block, uint, uint, error) {
	data, err := os.ReadFile(r.infoFilePath)
	if err != nil {
		zlog.Error("Read seat info file failed, err:", err)
		return nil, nil, 0, 0, err
	}
	var jsonSeat JsonSeat
	err = json.Unmarshal(data, &jsonSeat)
	if err != nil {
		zlog.Error("Unmarshal seat info failed, err:", err)
		return nil, nil, 0, 0, err
	}

	seats := make([][]*seat.Seat, 0)
	block := make([]*seat.Block, 0)

	// Make seats
	tmpSeats := make(map[uint]uint)
	for _, v := range jsonSeat.SeatInfo {
		start, end, err := getStartAndEnd(v.Row)
		if err != nil {
			zlog.Error("Get start and end failed, err:", err)
			return nil, nil, 0, 0, err
		}
		for i := start; i <= end; i++ {
			tmpSeats[i] = v.SeatNum
		}
	}
	for i := uint(0); i < jsonSeat.MaxRow; i++ {
		seats = append(seats, make([]*seat.Seat, 0))
		for j := uint(0); j < tmpSeats[i+1]; j++ {
			seats[i] = append(seats[i], &seat.Seat{
				Mutex:  sync.Mutex{},
				Row:    i + 1,
				Column: j + 1,
				Status: seat.Available,
			})
		}
	}
	// Make block
	for _, v := range jsonSeat.BlockInfo {
		start, end, err := getStartAndEnd(v.Row)
		if err != nil {
			zlog.Error("Get start and end failed, err:", err)
			return nil, nil, 0, 0, err
		}
		for i := start; i <= end; i++ {
			block = append(block, &seat.Block{
				Row:       i,
				Col:       v.Col,
				Direction: seat.BlockDirection(v.Direction),
				BlockNum:  v.BlockNum,
			})
		}
	}
	return seats, block, jsonSeat.MaxRow, jsonSeat.MaxCol, nil
}

type TestSeatInfoReader struct{}

func (r *TestSeatInfoReader) Init(string) error {
	return nil
}

func (r *TestSeatInfoReader) Read() ([][]*seat.Seat, []*seat.Block, uint, uint, error) {
	seats := make([][]*seat.Seat, 0)
	block := make([]*seat.Block, 0)
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
	return seats, block, 10, 10, nil
}
