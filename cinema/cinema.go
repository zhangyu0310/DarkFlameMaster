package cinema

import (
	"DarkFlameMaster/config"
	"DarkFlameMaster/seat"
	"errors"
	zlog "github.com/zhangyu0310/zlogger"
)

var cinema *Cinema

var (
	ErrInvalidRowOrCol = errors.New("invalid row or column")
)

func init() {
	cinema = &Cinema{}
}

type Cinema struct {
	maxRow uint
	maxCol uint

	Seats [][]*seat.Seat
}

func Init() error {
	cfg := config.GetGlobalConfig()
	var reader Reader
	switch cfg.SeatFileType {
	case config.JsonType:
		reader = &JsonTypeSeatInfoReader{}
	case config.TestType:
		reader = &TestSeatInfoReader{}
	default:
		reader = &TestSeatInfoReader{}
	}
	err := reader.Init(cfg.SeatFile)
	if err != nil {
		zlog.Error("Init seat file reader failed, err:", err)
		return err
	}
	seats, maxRow, maxCol, err := reader.Read()
	if err != nil {
		zlog.ErrorF("Read seat info [%s] failed, err: [%s]", cfg.SeatFile, err)
		return err
	}
	cinema.init(maxRow, maxCol, seats)
	return nil
}

func (c *Cinema) init(maxRow, maxCol uint, seats [][]*seat.Seat) {
	c.Seats = seats
	c.maxRow = maxRow
	c.maxCol = maxCol
}

func (c *Cinema) checkRowColValid(row, col uint) error {
	if c.maxRow < row || c.maxCol < col {
		return ErrInvalidRowOrCol
	}
	return nil
}

func (c *Cinema) GetSeat(row, col uint) *seat.Seat {
	return c.Seats[row-1][col-1]
}

func GetSeat(row, col uint) (*seat.Seat, error) {
	err := cinema.checkRowColValid(row, col)
	if err != nil {
		return nil, err
	}
	return cinema.GetSeat(row, col), nil
}

func AssociateSeats(seatInfo []*seat.Seat) []*seat.Seat {
	result := make([]*seat.Seat, 0, len(seatInfo))
	for _, s := range seatInfo {
		result = append(result, cinema.GetSeat(s.Row, s.Column))
	}
	return result
}

func GetSeatMap() ([]*seat.Seat, uint, uint) {
	maxRow := cinema.maxRow
	maxCol := cinema.maxCol
	seatInfo := make([]*seat.Seat, 0)

	for i := uint(0); i < maxRow; i++ {
		for j := uint(0); j < maxCol; j++ {
			seatInfo = append(seatInfo, cinema.Seats[i][j])
		}
	}

	return seatInfo, maxRow, maxCol
}
