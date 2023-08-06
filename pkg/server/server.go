package server

import (
	"fmt"
	"github.com/gucooing/gomcbe/pkg/config"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func Server() {
	var stdin io.WriteCloser
	var cmd *exec.Cmd
	var stdout io.ReadCloser
	var wg sync.WaitGroup

	cmd = exec.Command(config.GetConfig().ServerPath)
	stdin, _ = cmd.StdinPipe()
	stdout, _ = cmd.StdoutPipe()
	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, stdout)
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Println("无法启动游戏服务器:", err)
		return
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			fmt.Println("服务器退出原因:", err)
		} else {
			fmt.Println("服务器停止运行")
		}
	}()

	cmd1 := make(chan string)

	go Cmdenter(cmd1)
	for {
		select {
		case input := <-cmd1:
			fmt.Print(">> ")
			fmt.Print(input)
			input = strings.TrimSpace(input)
			if input == "stop" {
				break
			} else {
				_, err := io.WriteString(stdin, input+"\n")
				if err != nil {
					fmt.Println("无法发送命令原因:", err)
					break
				}
			}
		default:
			//fmt.Println("无输入:")
		}
	}
	err = stdin.Close()
	if err != nil {
		fmt.Println("无法关闭输入:", err)
	}
}
