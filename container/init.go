package container

import (
	"syscall"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"io/ioutil"
	"strings"
	"fmt"
)

func RunContainerInitProcess() error {

	//一直阻塞等待管道传递过来的命令
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, "")
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	//调用exec.LookPath可以在系统PATH内寻找命令的绝对路径 那么/bin/ls 就可以写为ls
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}


func readUserCommand() []string {
	//index为3的文件描述符,也就是传递进来的管道一端
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
