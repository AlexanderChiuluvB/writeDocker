package main

import (
	"./container"
	"github.com/Sirupsen/logrus"
	"os"
	"strings"
)

func Run(tty bool, commandArray []string){

	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("创建新的进程失败")
	}

	if err := parent.Start();err != nil {
		logrus.Error(err)
	}

	sendInitCommand(commandArray, writePipe)
	parent.Wait()
}

func sendInitCommand(commandArray []string, writePipe *os.File) {

	command := strings.Join(commandArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
