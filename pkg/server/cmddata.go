package server

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// 监听控制台输入，并将结果传入服务器
func Cmdenter(stdin io.WriteCloser) {
	defer stdin.Close()
	scanner := bufio.NewScanner(os.Stdin)
	//fmt.Printf(scanner)
	for scanner.Scan() {
		command := scanner.Text()
		_, err := fmt.Fprintln(stdin, command)
		if err != nil {
			fmt.Printf("无法输入: %v\n", err)
			break
		}
	}
}
