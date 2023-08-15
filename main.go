package main

import (
	"DarkFlameMaster/cinema"
	"DarkFlameMaster/config"
	"DarkFlameMaster/customer"
	"DarkFlameMaster/ticket/mgr"
	"DarkFlameMaster/web"
	zlog "github.com/zhangyu0310/zlogger"
	"os"
	"os/signal"
	"syscall"
)

func initConfig(cfg *config.Config) {
	cfg.SeatFile = "./data/seat.json"
	cfg.DbPath = "./run"
	_ = zlog.New("./run", "zlogger", true, zlog.LogLevelDebug)
}

func main() {
	config.InitializeConfig(initConfig)

	err := cinema.Init()
	if err != nil {
		zlog.Fatal("Init cinema failed, err:", err)
	}

	err = customer.Init()
	if err != nil {
		zlog.Fatal("Init customer info failed, err:", err)
	}

	err = mgr.Init()
	if err != nil {
		zlog.Fatal("Init ticket manager failed, err:", err)
	}

	srv := web.RunWebServer()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zlog.Info("Shutdown Server ...")

	web.ShutdownWebServer(srv)
}
