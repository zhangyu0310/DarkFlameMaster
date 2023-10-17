package ticket

import (
	"fmt"
	"time"
)

type Ticket struct {
	ID            string    `json:"id"`
	CustomerProof string    `json:"customerProof"`
	Row           uint      `json:"row"`
	Column        uint      `json:"column"`
	CreateTime    time.Time `json:"createTime"`
	Additional    string    `json:"additional"`
}

func NewTicket(proof string, row, column uint, additional string) *Ticket {
	return &Ticket{
		ID:            fmt.Sprintf("%s-%d-%d", proof, row, column),
		CustomerProof: proof,
		Row:           row,
		Column:        column,
		CreateTime:    time.Now(),
		Additional:    additional,
	}
}
