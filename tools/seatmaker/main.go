package main

import (
	"context"
	"github.com/gin-gonic/gin"
	zlog "github.com/zhangyu0310/zlogger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO: 实现
func main() {
	srv := RunWebServer()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zlog.Info("Shutdown Server ...")

	ShutdownWebServer(srv)
}

func RunWebServer() *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("view/*.html")

	router.Static("/script", "view/script")
	router.Static("/asset", "view/asset")
	router.GET("/", BaseSetting)
	router.POST("/seat_map", SeatMap)

	srv := &http.Server{Addr: ":9281", Handler: router}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.FatalF("listen: %s\n", err)
		}
	}()

	return srv
}

func ShutdownWebServer(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zlog.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		zlog.Warn("timeout of 5 seconds.")
	}
	zlog.Info("Server exiting")
}

func BaseSetting(context *gin.Context) {
	context.HTML(http.StatusOK, "base_setting.html", gin.H{})
}

func SeatMap(context *gin.Context) {
	maxRow := context.PostForm("maxRow")
	maxCol := context.PostForm("maxCol")
	order := context.PostForm("order")

	context.HTML(http.StatusOK, "seat_map.html", gin.H{
		"maxRow": maxRow,
		"maxCol": maxCol,
		"order":  order,
	})
}
