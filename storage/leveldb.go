package storage

import (
	"DarkFlameMaster/ticket"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	zlog "github.com/zhangyu0310/zlogger"
)

type TicketLevelDB struct {
	db *leveldb.DB
}

func (s *TicketLevelDB) Init(dbPath string, option ...interface{}) (err error) {
	var o *opt.Options
	if option != nil {
		o = option[0].(*opt.Options)
	}
	s.db, err = leveldb.OpenFile(dbPath, o)
	return
}

func (s *TicketLevelDB) ReadAll() ([]*ticket.Ticket, error) {
	tickets := make([]*ticket.Ticket, 0, 128)
	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		zlog.DebugF("Read all tickets key [%s], value [%s].",
			string(key), string(value))
		var t ticket.Ticket
		err := json.Unmarshal(value, &t)
		if err != nil {
			zlog.Error("Unmarshal ticket info [%s] from db failed, err:",
				string(value), err)
			return nil, err
		}
		tickets = append(tickets, &t)
	}
	iter.Release()
	return tickets, nil
}

func (s *TicketLevelDB) WriteTickets(tickets []*ticket.Ticket) error {
	batch := leveldb.Batch{}
	for _, t := range tickets {
		data, err := json.Marshal(t)
		if err != nil {
			zlog.Error("Marshal ticket failed, err:", err)
			return err
		}
		batch.Put([]byte(t.ID), data)
	}
	err := s.db.Write(&batch, &opt.WriteOptions{Sync: true})
	if err != nil {
		zlog.Error("Write ticket to db failed, err:", err)
		return err
	}
	return nil
}
