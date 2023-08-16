package cinema

import (
	"DarkFlameMaster/seat"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestJsonTypeSeatInfoReader_Read(t *testing.T) {
	r := &JsonTypeSeatInfoReader{}
	err := r.Init("./jsonTypeSeatInfoTest.json")
	assert.Nil(t, err)
	seats, block, maxRow, maxCol, err := r.Read()
	assert.Nil(t, err)
	assert.Equal(t, uint(15), maxRow)
	assert.Equal(t, uint(27), maxCol)
	// check seats
	for i := 0; i < 4; i++ {
		assert.Equal(t, 23, len(seats[i]))
	}
	for i := 4; i < 11; i++ {
		assert.Equal(t, 18, len(seats[i]))
	}
	for i := 11; i < 14; i++ {
		assert.Equal(t, 22, len(seats[i]))
	}
	for i := 14; i < 15; i++ {
		assert.Equal(t, 27, len(seats[i]))
	}
	// check block
	sort.Slice(block, func(i, j int) bool {
		if block[i].Row < block[j].Row {
			return true
		} else if block[i].Row > block[j].Row {
			return false
		}
		return block[i].Col < block[j].Col
	})

	for i := 0; i < 8; i++ {
		if i%2 == 0 {
			assert.Equal(t, uint(1), block[i].Col)
			assert.Equal(t, uint(2), block[i].BlockNum)
			assert.Equal(t, seat.DirectionFront, block[i].Direction)
		} else {
			assert.Equal(t, uint(23), block[i].Col)
			assert.Equal(t, uint(2), block[i].BlockNum)
			assert.Equal(t, seat.DirectionBack, block[i].Direction)
		}
	}
	for i := 8; i < 22; i++ {
		if i%2 == 0 {
			assert.Equal(t, uint(1), block[i].Col)
			assert.Equal(t, uint(5), block[i].BlockNum)
			assert.Equal(t, seat.DirectionFront, block[i].Direction)
		} else {
			assert.Equal(t, uint(18), block[i].Col)
			assert.Equal(t, uint(5), block[i].BlockNum)
			assert.Equal(t, seat.DirectionBack, block[i].Direction)
		}
	}
	for i := 22; i < 28; i++ {
		if i%2 == 0 {
			assert.Equal(t, uint(3), block[i].Col)
			assert.Equal(t, uint(3), block[i].BlockNum)
			assert.Equal(t, seat.DirectionFront, block[i].Direction)
		} else {
			assert.Equal(t, uint(20), block[i].Col)
			assert.Equal(t, uint(3), block[i].BlockNum)
			assert.Equal(t, seat.DirectionBack, block[i].Direction)
		}
	}
}
