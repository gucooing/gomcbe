package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
		for {
			programPath := config.GetConfig().ServerPath
			args := []string{config.GetConfig().Args}

			cmd, stdin, err := server.Server(programPath, args)
			if err != nil {
				fmt.Printf("服务器启动失败: %v\n", err)
				return
			}
			defer stdin.Close()
			go server.Cmdenter(stdin)
			fmt.Println("\nMC 启动！")
			err = cmd.Wait()
			break
		}
		//函数退出，触发进程守护重启服务器
		fmt.Println("2秒后重启服务器...")
		time.Sleep(2 * time.Second)
	}
}
