package web

import (
	"DarkFlameMaster/cinema"
	"DarkFlameMaster/customer"
	"DarkFlameMaster/seat"
	"DarkFlameMaster/ticket/tkmgr"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	zlog "github.com/zhangyu0310/zlogger"
	"net/http"
	"time"
)

var (
	GinRunMode string
)

func RunWebServer() *http.Server {
	if GinRunMode != gin.ReleaseMode {
		GinRunMode = ""
	}
	gin.SetMode(GinRunMode)
	router := gin.Default()
	router.LoadHTMLGlob("view/*.html")

	router.Static("/script", "view/script")
	router.Static("/asset", "view/asset")
	router.GET("/", Proof)
	router.GET("/choose_seat", ChooseSeat)
	router.POST("/choose_seat", ChooseResult)
	router.POST("/check_tickets", CheckTickets)

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

type SeatWrapper struct {
	Row    uint        `json:"row"`
	Column uint        `json:"col"`
	Status seat.Status `json:"status"`
	IsMine bool        `json:"isMine"`
}

type SendMsg struct {
	Proof        string         `json:"proof"`
	Msg          string         `json:"msg"`
	CanChooseNum uint           `json:"canChooseNum"`
	MaxRow       uint           `json:"maxRow"`
	MaxCol       uint           `json:"maxCol"`
	Order        string         `json:"order"`
	SeatInfo     []*SeatWrapper `json:"seatInfo"`
	BlockInfo    []*seat.Block  `json:"blockInfo"`
}

type ReceiveMsg struct {
	Proof    string       `json:"proof"`
	SeatInfo []*seat.Seat `json:"seatInfo"`
}

func Proof(context *gin.Context) {
	context.HTML(http.StatusOK, "proof.html",
		gin.H{
			"title":     "Dark Flame Master - Choose Seat",
			"proofName": "QQ号",
		})
}

func sendDataToWeb(context *gin.Context, c *customer.Customer, errMsg string) {
	// 打包发送数据
	si, bl, maxRow, maxCol, order := cinema.GetSeatMap()
	seatWrapper := make([]*SeatWrapper, 0, len(si))
	for _, s := range si {
		sw := &SeatWrapper{
			Row:    s.Row,
			Column: s.Column,
			Status: s.Status,
			IsMine: false,
		}
		if s.Ticket != nil && s.Ticket.CustomerProof == c.Proof {
			sw.IsMine = true
		}
		seatWrapper = append(seatWrapper, sw)
	}
	sendData := SendMsg{
		Proof:        c.Proof,
		Msg:          errMsg,
		CanChooseNum: c.RemainTicketNumber(),
		MaxRow:       maxRow,
		MaxCol:       maxCol,
		Order:        order,
		SeatInfo:     seatWrapper,
		BlockInfo:    bl,
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
		errMsg = "选座数据解析失败！需要重新登入。"
		context.HTML(http.StatusOK, "proof.html",
			gin.H{
				"title": "Dark Flame Master - Choose Seat",
				"error": errMsg,
			})
	} else {
		cus, err = customer.GetCustomer(reMsg.Proof)
		if err != nil {
			zlog.Fatal("GetCustomer failed when choose result, err:", err)
		}
		seats := cinema.AssociateSeats(reMsg.SeatInfo)
		tk, err := tkmgr.MakeTickets(cus, seats)
		if err != nil {
			zlog.Error("Make tickets failed, err:", err)
			errMsg = "选座失败！"
		}
		zlog.DebugF("Tickets: %v", tk)
		sendDataToWeb(context, cus, errMsg)
	}
}

func CheckTickets(context *gin.Context) {
	proof := context.PostForm("check")
	t := tkmgr.CheckTickets(proof)
	msg := ""
	if len(t) == 0 {
		msg = "当前用户没有选座记录！"
	} else {
		for _, v := range t {
			msg += fmt.Sprintf("第%d排，第%d座，选座时间:%s\n",
				v.Row, v.Column, v.CreateTime.Format("2006-01-02 15:04:05"))
		}
	}
	context.String(http.StatusOK, msg)
	return
}
