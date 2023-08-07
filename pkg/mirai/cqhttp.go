package mirai

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/gucooing/gomcbe/pkg/config"
	"github.com/gucooing/gomcbe/pkg/server"
)

var conn *websocket.Conn = nil

// Reqws 函数用于建立与 cqhttp 的 WebSocket 连接
func Reqws(ctx context.Context) {
	// 创建 WebSocket 连接
	var err error
	serverURL := config.GetConfig().CqhttpWsurl
	conn, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("收到重启信号，重置WS连接")
			// 关闭 WebSocket 连接
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			_ = reswsdata(string(message))
		}
	}
}

type Reqdata struct {
	MessageType string  `json:"message_type"`
	Sender      *Sender `json:"sender"`
	Message     string  `json:"message"`
}

type Sender struct {
	Nickname string `json:"nickname"`
	UserId   int    `json:"user_id"`
}

func reswsdata(message string) string {
	var reqdata Reqdata
	err := json.Unmarshal([]byte(message), &reqdata)
	if err != nil {
		fmt.Println("解析 JSON 出错:", err)
		return ""
	}
	fmt.Printf("ws接收数据: %v\n", message)
	if reqdata.Sender == nil {
		return ""
	}
	if reqdata.Sender.UserId == config.GetConfig().QqAdmin {
		server.Sender(reqdata.Message)
	}
	return ""
}

// SendWSMessage 定义发送函数
func SendWSMessage(msg interface{}) error {
	var err error
	serverURL := config.GetConfig().CqhttpWsurl
	conn, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	// 发送消息
	err = conn.WriteJSON(msg)
	if err != nil {
		return err
	}
	return nil
}

type Params struct {
	MessageType string `json:"message_type"`
	UserId      int    `json:"user_id"`
	Message     string `json:"message"`
	AutoEscape  bool   `json:"auto_escape"`
}

type Rsqdata struct {
	Action string  `json:"action"`
	Params *Params `json:"params"`
}

// SendWSMessagesi 定义私聊发送函数
func SendWSMessagesi(msg string) {
	fmt.Println(config.GetConfig().QqAdmin)
	rsqdata := Rsqdata{
		Action: "send_msg",
		Params: &Params{
			MessageType: "private",
			UserId:      config.GetConfig().QqAdmin,
			Message:     msg,
			AutoEscape:  false,
		},
	}
	// 发送消息
	err := SendWSMessage(rsqdata)
	if err != nil {
		return
	}
}
