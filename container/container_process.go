package container

import (
	"os"
	"os/exec"
	"syscall"
)


func NewParentProcess(tty bool, command string) *exec.Cmd {

	args := []string{"init", command}
	// /proc/self指当前运行进程自己的环境
	// exe就是自己调用了自己,用这种方式来进行初始化
	cmd := exec.Command("/proc/self/exe", args...)
	//fork出来一个新的进程,使用namespace隔离新创建的进程与外部环境
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	return cmd
}

