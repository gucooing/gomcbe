package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/gucooing/gomcbe/pkg/config"
	"github.com/gucooing/gomcbe/pkg/mirai"
	"github.com/gucooing/gomcbe/pkg/server"
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
	for {
		//启动服务器
		// 创建一个带有取消机制的上下文（context）
		ctx, cancel := context.WithCancel(context.Background())
		programPath := config.GetConfig().ServerPath
		cmd, stdin, err := server.Server(programPath)
		if err != nil {
			fmt.Printf("服务器启动失败: %v\n", err)
			return
		}
		defer func(cmd *exec.Cmd) {
			// 停止服务器进程
			err := cmd.Process.Kill()
			if err != nil {
				fmt.Println("无法终止服务器进程:", err)
			}
		}(cmd)
		mirai.SendWSMessagesi("服务器 启动！")
		fmt.Println("\nServer 启动！")
		//启动控制台监控
		go server.Command(stdin, ctx)
		//启动第三方命令传入通道
		go server.Commands(stdin, ctx)
		//连接ws
		go mirai.Reqws(ctx)
		err = cmd.Wait()
		mirai.SendWSMessagesi("服务器正在关闭")
		//函数退出，触发进程守护重启服务器
		fmt.Println("发送停止信息...")
		cancel()
		fmt.Println("2秒后重启服务器...")
		time.Sleep(2 * time.Second)
	}
}
