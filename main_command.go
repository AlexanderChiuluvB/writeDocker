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
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
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
		cmd := ctx.Args().Get(0)
		logrus.Info("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
