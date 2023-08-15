package web

import (
	"DarkFlameMaster/cinema"
	"DarkFlameMaster/customer"
	"DarkFlameMaster/seat"
	"DarkFlameMaster/ticket/mgr"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	zlog "github.com/zhangyu0310/zlogger"
	"net/http"
	"time"
)

func RunWebServer() *http.Server {
	router := gin.Default()
	router.LoadHTMLGlob("view/*.html")

	router.Static("/script", "view/script")
	router.Static("/asset", "view/asset")
	router.GET("/", Proof)
	router.GET("/choose_seat", ChooseSeat)
	router.POST("/choose_seat", ChooseResult)
	// TODO: 查票相关代码

	srv := &http.Server{Addr: ":718", Handler: router}
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

type SendMsg struct {
	Proof        string       `json:"proof"`
	Msg          string       `json:"msg"`
	CanChooseNum uint         `json:"canChooseNum"`
	MaxRow       uint         `json:"maxRow"`
	MaxCol       uint         `json:"maxCol"`
	SeatInfo     []*seat.Seat `json:"seatInfo"`
}

type ReceiveMsg struct {
	Proof    string       `json:"proof"`
	SeatInfo []*seat.Seat `json:"seatInfo"`
}

func Proof(context *gin.Context) {
	context.HTML(http.StatusOK, "proof.html",
		gin.H{
			"title": "Dark Flame Master - Choose Seat",
		})
}

func sendDataToWeb(context *gin.Context, c *customer.Customer, errMsg string) {
	// 打包发送数据
	si, maxRow, maxCol := cinema.GetSeatMap()
	sendData := SendMsg{
		Proof:        c.Proof,
		Msg:          errMsg,
		CanChooseNum: c.RemainTicketNumber(),
		MaxRow:       maxRow,
		MaxCol:       maxCol,
		SeatInfo:     si,
	}
	data, err := json.Marshal(sendData)
	if err != nil {
		zlog.Error("Marshal seat info data failed. err:", err)
		context.HTML(http.StatusOK, "choose_seat.html",
			gin.H{
				"title":    "Dark Flame Master - Choose Seat",
				"seatData": "{\"msg\":\"Server have fatal error.\"}",
			})
	} else {
		context.HTML(http.StatusOK, "choose_seat.html",
			gin.H{
				"title":    "Dark Flame Master - Choose Seat",
				"seatData": string(data),
			})
	}
}

func ChooseSeat(context *gin.Context) {
	proof := context.Query("proof")
	cus, err := customer.GetCustomer(proof)
	errMsg := ""
	if err != nil {
		zlog.DebugF("GetCustomer failed. Proof [%s] err: %s", proof, err)
		errMsg = "交易单号不存在！请检查后重新选座。"
	} else {
		if !customer.CanChooseSeat(cus) {
			errMsg = "当前订单暂时还无法进行选座！"
		}
	}
	if errMsg != "" {
		context.HTML(http.StatusOK, "proof.html",
			gin.H{
				"title": "Dark Flame Master - Choose Seat",
				"error": errMsg,
			})
	} else {
		sendDataToWeb(context, cus, errMsg)
	}
}

func ChooseResult(context *gin.Context) {
	data := context.PostForm("chooseData")
	errMsg := "选座成功！"
	var cus *customer.Customer
	reMsg := ReceiveMsg{}
	err := json.Unmarshal([]byte(data), &reMsg)
	if err != nil {
		zlog.Error("Unmarshal seat choose data failed. err:", err)
		errMsg = "选座数据解析失败！"
	} else {
		cus, err = customer.GetCustomer(reMsg.Proof)
		if err != nil {
			zlog.Fatal("GetCustomer failed when choose result, err:", err)
		}
		seats := cinema.AssociateSeats(reMsg.SeatInfo)
		tk, err := mgr.MakeTickets(cus, seats)
		if err != nil {
			zlog.Error("Make tickets failed, err:", err)
			errMsg = "选座失败！"
		}
		zlog.DebugF("Tickets: %v", tk)
	}
	sendDataToWeb(context, cus, errMsg)
}
