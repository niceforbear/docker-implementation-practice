package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
	"os"
)

const usage = `Docker implementation practice.`

func main() {
	//setupUtsNs()
	//setupIpcNs()
	//setupPidNs()
	//setupMountNs()
	//setupUserNs()
	//setupNetNs()
	//setupCgroupMemory()

	app := cli.App{}
	app.Name = "mydocker"
	app.Usage = usage

	var commands []*cli.Command
	commands = append(commands, &initCommand)
	commands = append(commands, &runCommand)
	commands = append(commands, &commitCommand)
	app.Commands = commands

	// 初始化日志配置
	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})

		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}


