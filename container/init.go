package container

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

/**
这里的init是在容器内进行的,这是本容器执行的第一个进程.
使用mount先去挂载proc文件系统


初始化容器内容,挂载proc文件系统,运行用户指定程序
 */
func RunContainerInitProcess(command string, args []string) error {

	mount()

	path, err := exec.LookPath(command)
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}

	argv := []string{command}
	log.Infof("Find path %s", path)
	if err := syscall.Exec(command, argv, os.Environ());err!=nil{
		/*
		 */
		log.Errorf(err.Error())
	}
	return nil
}

func mount() {

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID|
		syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc",uintptr(defaultMountFlags), "")
}


