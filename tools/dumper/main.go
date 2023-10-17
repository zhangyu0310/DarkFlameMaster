package main

import (
	"DarkFlameMaster/tools/dumper/dump"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

var (
	dbPath = *flag.String("db-path", "./run/db", "db path")
	dbType = *flag.String("db-type", "leveldb", "db type")
)

func main() {
	var dumper dump.Dump
	switch dbType {
	case "leveldb":
		dumper = &dump.LevelDBDump{}
	default:
		panic("not support db type")
	}
	tickets, err := dumper.Dump(dbPath)
	if err != nil {
		panic(err)
	}
	sort.Slice(tickets, func(i, j int) bool {
		return tickets[i].CreateTime.Before(tickets[j].CreateTime)
	})
	fileName := fmt.Sprintf("dump_tickets_%s.csv",
		time.Now().Format("20060102150405"))
	// create dump file and write data
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	err = dump.ToCSV(file, tickets)
	if err != nil {
		panic(err)
	}
	_ = file.Close()
}
