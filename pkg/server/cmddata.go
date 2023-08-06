package server

import (
	"bufio"
	"fmt"
	"os"
)

// 监听控制台输入，并将结果传入服务器
func Cmdenter(cm chan<- string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		cm <- fmt.Sprintf(input)
	}
}
