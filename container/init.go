package container

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

/**
这里的init是在容器内进行的,这是本容器执行的第一个进程.
使用mount先去挂载proc文件系统
初始化容器内容,挂载proc文件系,运行用户指定程序
 */
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) ==0 {
		return fmt.Errorf("参数为无")
	}

	mount()

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)
	//系统调用实现了完成初始化的操作,并且将用户进程运行起来
	if err := syscall.Exec(path, cmdArray[0:], os.Environ());err!=nil{
		log.Errorf(err.Error())
	}
	return nil
}

func mount() {
	pwd,err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return
	}
	log.Infof("当前挂载位置为",pwd)

	//pivotROOT
	if err := pivotRoot(pwd);err != nil {
		log.Errorf("error occured in pivot function is: %v", err)
		return
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID|
		syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc",uintptr(defaultMountFlags), "")
	syscall.Mount("dev","/dev","tmpfs",syscall.MS_NOSUID|syscall.MS_STRICTATIME,
		"mode=755")
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer pipe.Close()
	msg,err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("读取管道初始化失败 ",err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}


func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if _, err := os.Stat(pivotDir); os.IsNotExist(err) {
		if err = os.Mkdir(pivotDir, 0777); err != nil {
			return err
		}
	}

	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("new Get current location error %v", err)
	}
	println("new Current location is ", pwd)
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	return os.Remove(pivotDir)
}
