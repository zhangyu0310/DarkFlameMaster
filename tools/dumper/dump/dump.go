package dump

import (
	"DarkFlameMaster/ticket"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
)

type Dump interface {
	Dump(string) ([]*ticket.Ticket, error)
}

type LevelDBDump struct{}

func (ld *LevelDBDump) Dump(path string) ([]*ticket.Ticket, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	iter := db.NewIterator(nil, nil)
	defer iter.Release()
	tickets := make([]*ticket.Ticket, 0, 128)
	for iter.Next() {
		value := iter.Value()
		var t ticket.Ticket
		err := json.Unmarshal(value, &t)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, &t)
	}
	return tickets, nil
}
