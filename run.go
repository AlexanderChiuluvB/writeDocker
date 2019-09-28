package main

import (
	"./container"
	"github.com/Sirupsen/logrus"
	"os"
)

func Run(tty bool, command string){

	parent:= container.NewParentProcess(tty, command)
	if parent == nil {
		logrus.Errorf("创建新的进程失败")
	}

	if err := parent.Start();err != nil {
		/*
		会调用前面创建的command的进程
		1.首先clone出来一个namespace隔离的进程
		2.在子进程中,调用/proc/self/exe 也就是自己,发送init参数
		3.调用我们initCommand方法,初始化容器的资源
		 */
		logrus.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}

