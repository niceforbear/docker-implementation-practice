package main

import (
	"github.com/niceforbear/docker-implementation-practice/container"
	"github.com/sirupsen/logrus"
	"os"
)

// 作用：容器内进程调用自己
func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)

	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	err := parent.Wait()
	if err != nil {
		logrus.Error(err)
	}

	os.Exit(-1)
}