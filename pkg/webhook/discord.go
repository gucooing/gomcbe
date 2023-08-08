package webhook

import (
	"bytes"
	"fmt"
	"github.com/gucooing/gomcbe/pkg/config"
	"io/ioutil"
	"net/http"
)

func Discord(msg string) {
	url := config.GetConfig().DiscordWebhookUrl
	body := "{\"content\": \"" + msg + "\", \"username\": \"进程守护服务控制台发送\", \"avatar_url\": \"https://go.dev/images/favicon-gopher.png\"}"
	response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Println("discord webhook消息发送错误:", err)
	}
	b, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(b))
}
