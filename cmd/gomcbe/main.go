package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gucooing/gomcbe/pkg/config"
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
		programPath := config.GetConfig().ServerPath
		cmd, stdin, err := server.Server(programPath)
		if err != nil {
			fmt.Printf("服务器启动失败: %v\n", err)
			return
		}
		defer func(stdin io.WriteCloser) {
			err := stdin.Close()
			if err != nil {

			}
		}(stdin)
		//启动控制台监控
		go server.Command(stdin)
		//启动第三方命令传入通道
		go server.Commands(stdin)
		fmt.Println("\nServer 启动！")
		err = cmd.Wait()
		//函数退出，触发进程守护重启服务器
		fmt.Println("2秒后重启服务器...")
		time.Sleep(2 * time.Second)
	}
}
