package main

import (
	"fmt"
	"github.com/niceforbear/docker-implementation-practice/container"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

// 作用：解析参数
var runCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with ns and cgroups limit mydocker run -ti [cmd]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},

	/*
		1. 判断参数是否包含 cmd， 获取用户指定的 cmd。
		2. 调用 Run function 去准备启动容器。
	*/
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("missing container cmd")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command{
	Name: "init",
	Usage: "Init container process",
	/*
		1. 获取传递过来的 cmd 参数
		2. 执行容器初始化
	*/
	Action: func(context *cli.Context) error {
		log.Infof("init come on")

		cmd := context.Args().Get(0)

		log.Infof("cmd %v", cmd)

		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}

var runCommandV2 = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run -ti [command]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		&cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
	},
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args().Slice() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		volume := context.String("v")
		RunV5(tty, cmdArray, volume)
		return nil
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit a container into image",
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("Missing container name")
		}
		imageName := context.Args().Get(0)
		//commitContainer(containerName)
		commitContainer(imageName)
		return nil
	},
}