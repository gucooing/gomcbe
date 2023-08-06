package server

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func Server(path string) (*exec.Cmd, io.WriteCloser, error) {
	var cmd *exec.Cmd
	var err error

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c", "chcp 65001 &&", path)
	} else {
		cmd = exec.Command("sh", "-c", path)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("无法输入: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("启动失败: %v", err)
	}

	return cmd, stdin, nil
}
