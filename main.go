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

func initConfig(cfg *config.Config) {
	config.ReadConfig("./conf/configure.json", cfg)
	_ = zlog.New(cfg.LogPath, "zlogger", true, zlog.LogLevelDebug)
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
