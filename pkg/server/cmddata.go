package server

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Command 监听控制台输入，并将结果传入服务器
func Command(stdin io.WriteCloser) {
	defer func(stdin io.WriteCloser) {
		err := stdin.Close()
		if err != nil {

		}
	}(stdin)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		SendCommand(stdin, command)
	}
}

// Commands 接收第三方命令，传入服务器
func Commands(stdin io.WriteCloser) {
	defer func(stdin io.WriteCloser) {
		err := stdin.Close()
		if err != nil {

		}
	}(stdin)
	commandment := make(chan string)
	for {
		select {
		case command := <-commandment:
			SendCommand(stdin, command)
		}
	}
}

// Sender 第三方命令传入通道
func Sender(msg string) {
	commandment := make(chan string)
	commandment <- msg
}

// SendCommand 发送命令到服务器
func SendCommand(stdin io.WriteCloser, command string) {
	_, err := fmt.Fprintln(stdin, command)
	if err != nil {
		fmt.Printf("无法输入: %v\n", err)
	}
}
