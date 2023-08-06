package server

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

func Server(path string, args []string) (*exec.Cmd, io.WriteCloser, error) {
	cmd := exec.Command("cmd.exe")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("无法输入: %v", err)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: fmt.Sprintf("/c chcp 65001 && \"%s\"", path),
	}
	err = cmd.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("启动失败: %v", err)
	}
	return cmd, stdin, nil
}
