package main

import (
	"DarkFlameMaster/cinema"
	"DarkFlameMaster/config"
	"DarkFlameMaster/customer"
	"DarkFlameMaster/ticket/tkmgr"
	"DarkFlameMaster/web"
	zlog "github.com/zhangyu0310/zlogger"
	"os"
	"os/signal"
	"syscall"
)

// TODO: read config from file
func initConfig(cfg *config.Config) {
	cfg.DbType = config.LevelDB
	cfg.DbPath = "./run/db"
	cfg.SeatFileType = config.JsonType
	cfg.SeatFile = "./data/奥斯卡长安国际影城-5号ALPD激光厅.json"
	cfg.CustomerType = config.NoPay
	cfg.CustomerFile = "./data/customer.json"
	cfg.ChooseSeatStrategy = config.NoLimit
	_ = zlog.New("./run/log", "zlogger", true, zlog.LogLevelDebug)
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

	err = tkmgr.Init()
	if err != nil {
		zlog.Fatal("Init ticket manager failed, err:", err)
	}

	srv := web.RunWebServer()
	web.RunAdminServer()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zlog.Info("Shutdown Server ...")

	web.ShutdownWebServer(srv)
}
