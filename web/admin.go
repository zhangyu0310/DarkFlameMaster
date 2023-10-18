package web

import (
	"DarkFlameMaster/config"
	"DarkFlameMaster/seat"
	"DarkFlameMaster/ticket/tkmgr"
	"DarkFlameMaster/tools/dumper/dump"
	"encoding/json"
	"fmt"
	zlog "github.com/zhangyu0310/zlogger"
	"io"
	"net/http"
	"sort"
	"strings"
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
	// write csv data to response
	err = dump.ToCSV(w, tk)
	if err != nil {
		zlog.ErrorF("dump tickets failed, err: %s\n", err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}

type DeleteTicketsReq struct {
	Mode  string       `json:"mode"`
	Ids   []string     `json:"ids"`
	Seats []*seat.Seat `json:"seats"`
	Users []string     `json:"users"`
}

func deleteTickets(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		zlog.Error("read request body failed, err:", err)
		_, _ = w.Write([]byte(err.Error()))
	}
	req := &DeleteTicketsReq{}
	err = json.Unmarshal(data, req)
	if err != nil {
		zlog.Error("unmarshal request body failed, err:", err)
		_, _ = w.Write([]byte(err.Error()))
	}
	switch strings.ToLower(req.Mode) {
	case "id":
		err = tkmgr.DeleteTicketsByIds(req.Ids)
	case "seat":
		err = tkmgr.DeleteTicketsBySeats(req.Seats)
	case "user":
		err = tkmgr.DeleteTicketsByUsers(req.Users)
	default:
		zlog.Error("invalid mode:", req.Mode)
		_, _ = w.Write([]byte("invalid mode"))
	}
	if err != nil {
		zlog.Error("delete tickets failed, err:", err)
		_, _ = w.Write([]byte(err.Error()))
	} else {
		_, _ = w.Write([]byte("ok"))
	}
}

func RunAdminServer() {
	cfg := config.GetGlobalConfig()
	var addr string
	if cfg.AdminLocalOnly {
		addr = fmt.Sprintf("127.0.0.1:%d", cfg.AdminPort)
	} else {
		addr = fmt.Sprintf(":%d", cfg.AdminPort)
	}
	http.HandleFunc("/dump_tickets", dumpTickets)
	http.HandleFunc("/delete_tickets", deleteTickets)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil &&
			err != http.ErrServerClosed {
			zlog.FatalF("listen: %s\n", err)
		}
	}()
}
