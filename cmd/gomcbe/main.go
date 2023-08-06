package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gucooing/gomcbe/pkg/config"
)

func main() {
	// 启动读取配置
	err := config.LoadConfig()
	if err != nil {
		if err == config.FileNotExist {
			p, _ := json.MarshalIndent(config.DefaultConfig, "", "  ")
			fmt.Printf("找不到配置文件，这是默认配置:\n%s\n", p)
			fmt.Printf("\n您可以将其保存到名为“config.json”的文件中并再次运行该程序\n")
			fmt.Printf("Press 'Enter' to exit ...\n")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			os.Exit(0)
		} else {
			panic(err)
		}
	}

	// 启动守护服务
	reader := bufio.NewReader(os.Stdin)
	var cmd *exec.Cmd
	var stdin io.WriteCloser
	var stdout io.ReadCloser
	var wg sync.WaitGroup
	for {
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
		// 第一个循环：处理用户输入并将命令发送到游戏服务器
		for {
			fmt.Print(">> ")
			input, _ := reader.ReadString('\n')
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
		}
		err = stdin.Close()
		if err != nil {
			fmt.Println("无法关闭输入:", err)
		}
		cmd.Process.Kill()
		wg.Wait()

		// 第二个循环：等待用户命令，并执行相应操作
		for {
			fmt.Print("请输入命令 (start/stop): ")
			command, _ := reader.ReadString('\n')
			command = strings.TrimSpace(command)

			if command == "start" {
				break
			} else if command == "stop" {
				fmt.Println("正在停止服务器...")
				return
			} else {
				fmt.Println("无效的命令，请重新输入")
			}
		}

		fmt.Println("2秒后重启服务器...")
		time.Sleep(2 * time.Second)
	}
}
