package main

import (
	"DarkFlameMaster/cinema"
	"DarkFlameMaster/config"
	"DarkFlameMaster/customer"
	"DarkFlameMaster/serverinfo"
	"DarkFlameMaster/ticket/tkmgr"
	"DarkFlameMaster/web"
	"flag"
	"fmt"
	zlog "github.com/zhangyu0310/zlogger"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("config", "./conf/configure.json", "config path - json format")

func initConfig(cfg *config.Config) {
	config.ReadConfig(*configPath, cfg)
	_ = zlog.New(cfg.LogPath, "zlogger", true, zlog.LogLevelDebug)
}

func main() {
	serverinfo.InitInformation()
	v := flag.Bool("v", false, "show version")
	vv := flag.Bool("version", false, "show version")
	h := flag.Bool("h", false, "show help")
	hh := flag.Bool("help", false, "show help")
	flag.Parse()
	if *v || *vv {
		fmt.Println(serverinfo.Get().String())
		return
	}
	if *h || *hh {
		flag.Usage()
		return
	}

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

	web.GinRunMode = serverinfo.Get().BuildMode
	srv := web.RunWebServer()
	web.RunAdminServer()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zlog.Info("Shutdown Server ...")

	web.ShutdownWebServer(srv)
}
