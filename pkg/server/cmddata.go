package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
)

var Commandment = make(chan string)

// Command 监听控制台输入，并将结果传入服务器
func Command(stdin io.WriteCloser, ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	commandChan := make(chan string)
	// 开启一个goroutine监听控制台输入
	go func() {
		for scanner.Scan() {
			command := scanner.Text()
			commandChan <- command
		}
	}()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("收到停止信号，关闭Command")
			return
		case command := <-commandChan:
			SendCommand(stdin, command)
		}
	}
}

// Commands 接收第三方命令，传入服务器
func Commands(stdin io.WriteCloser, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("收到停止信号，关闭Commands")
			return
		case command := <-Commandment:
			fmt.Printf("QQ管理员命令: %v\n", command)
			SendCommand(stdin, command)
		}

	}
}

// Sender 第三方命令传入通道
func Sender(msg string) {
	Commandment <- msg
}

// SendCommand 发送命令到服务器
func SendCommand(stdin io.WriteCloser, command string) {
	_, err := fmt.Fprintln(stdin, command)
	if err != nil {
		fmt.Printf("无法输入: %v\n", err)
	}
}
