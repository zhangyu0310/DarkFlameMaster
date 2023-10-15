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
	*SeatData
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
	seatData, err := reader.Read()
	if err != nil {
		zlog.ErrorF("Read seat info [%s] failed, err: [%s]", cfg.SeatFile, err)
		return err
	}
	cinema.init(seatData)
	return nil
}

func (c *Cinema) init(data *SeatData) {
	c.SeatData = data
}

func (c *Cinema) checkRowColValid(row, col uint) error {
	if c.MaxRow < row || c.MaxCol < col {
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

func GetSeatMap() ([]*seat.Seat, []*seat.Block, uint, uint, string) {
	maxRow := cinema.MaxRow
	maxCol := cinema.MaxCol
	seatInfo := make([]*seat.Seat, 0)

	for i := uint(0); i < maxRow; i++ {
		for _, v := range cinema.Seats[i] {
			seatInfo = append(seatInfo, v)
		}
	}

	return seatInfo, cinema.Block, maxRow, maxCol, cinema.Order
}
