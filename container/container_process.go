package container

import (
	"os"
	"os/exec"
	"syscall"
	log "github.com/Sirupsen/logrus"
)


func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	/**
	外带着这个文件句柄取创建子进程
	一个进程默认有3个文件描述符,分别是标准输入,输出与错误
	这三个文件描述符是子进程创建的时候就默认附带着,所以外带的文件描述符会变成第四个
	 */
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}