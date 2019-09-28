package main


import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"./container"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: `my docker run -it [command], Create a docker with namespace and cgroup limits`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},

	/**
	执行run命令执行的真正函数
	获取参数,调用Run,准备启动容器
	**/
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")

		}
		var cmdArray []string
		for _, cmd := range context.Args() {
			cmdArray = append(cmdArray, cmd)
		}
		tty := context.Bool("ti")
		//资源限制
		Run(tty, cmdArray)
		return nil
	},
}

var initCommand = cli.Command{

	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",

	/**

	定义了initCommand具体操作

	*/
	Action: func(ctx *cli.Context) error {
		logrus.Info("init start")
		err := container.RunContainerInitProcess()
		return err
	},
}
