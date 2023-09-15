package test

import (
	g "InstantMessaging/app/global"
	"InstantMessaging/app/internal/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// PublishKey Redis中用于发布和订阅的通道名
const (
	PublishKey = "websocket"
)

// Publish 将消息发布到Redis指定的通道
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish ", msg)
	err = g.Rdb.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 从Redis订阅指定的通道，并返回接收到的消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := g.Rdb.Subscribe(ctx, channel)
	fmt.Println("Subscribe ", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe 。。。。", msg.Payload)
	return msg.Payload, err
}

// 设置UpGrader
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SendMsg 处理HTTP请求并将其升级为WebSocket连接
func SendMsg(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(c, ws)
}

// MsgHandler 实际的消息处理逻辑
func MsgHandler(c *gin.Context, ws *websocket.Conn) {
	for {
		msg, err := Subscribe(c, PublishKey)
		if err != nil {
			fmt.Println(" MsgHandler 发送失败", err)
		}

		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	service.Chat(c.Writer, c.Request)
}
