package main

import (
	"./container"
	log "github.com/Sirupsen/logrus"
	"./cgroups"
	"./cgroups/subsystems"
	"os"
	"strings"
)

func Run(tty bool, commandArray[] string, volume string, res *subsystems.ResourceConfig){

	parent, writePipe := container.NewParentProcess(tty,volume )

	if parent == nil {
		log.Errorf("创建新的进程失败")
	}

	if err := parent.Start();err != nil {
		/*
		会调用前面创建的command的进程
		1.首先clone出来一个namespace隔离的进程
		2.在子进程中,调用/proc/self/exe 也就是自己,发送init参数
		3.调用我们initCommand方法,初始化容器的资源
		 */
		log.Error(err)
	}

	//初始化subsystem实例
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	//遍历调用各个subsystem实例的set方法,创建与配置不同subsystem挂载的cgroup
	cgroupManager.Set(res)
	//把当前容器进程ID加入到那些cgroup
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(commandArray, writePipe)
	parent.Wait()
	mntURL := "/opt/test2/mnt/"
	rootURL := "/opt/test2/"
	container.DeleteWorkSpace(rootURL, mntURL, volume)
	os.Exit(-1)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
