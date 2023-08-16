package web

import (
	"DarkFlameMaster/ticket/tkmgr"
	"encoding/csv"
	"fmt"
	zlog "github.com/zhangyu0310/zlogger"
	"net/http"
	"os"
	"sort"
	"time"
)

func dumpTickets(w http.ResponseWriter, _ *http.Request) {
	// dump tickets from database
	tk, err := tkmgr.Dump()
	if err != nil {
		zlog.ErrorF("dump tickets failed, err: %s\n", err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	sort.Slice(tk, func(i, j int) bool {
		return tk[i].CreateTime.Before(tk[j].CreateTime)
	})
	// create dump file and write data
	fileName := fmt.Sprintf("dump_tickets_%s.csv",
		time.Now().Format("20060102150405"))
	file, err := os.Create(fileName)
	if err != nil {
		zlog.Error("create file failed, err:", err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	writer := csv.NewWriter(file)
	_ = writer.Write([]string{"ID", "CustomerProof", "Row", "Column", "CreateTime"})
	for _, t := range tk {
		_ = writer.Write([]string{
			t.ID,
			t.CustomerProof,
			fmt.Sprintf("%d", t.Row),
			fmt.Sprintf("%d", t.Column),
			t.CreateTime.Format("2006-01-02 15:04:05")})
	}
	_, _ = w.Write([]byte("ok"))
}

func deleteTickets(_ http.ResponseWriter, _ *http.Request) {

}

func RunAdminServer() {
	http.HandleFunc("/dump_tickets", dumpTickets)
	http.HandleFunc("/delete_tickets", deleteTickets)
	go func() {
		if err := http.ListenAndServe("127.0.0.1:1219", nil); err != nil &&
			err != http.ErrServerClosed {
			zlog.FatalF("listen: %s\n", err)
		}
	}()
}
