package main

import (
	"github.com/niceforbear/docker-implementation-practice/cgroups"
	"github.com/niceforbear/docker-implementation-practice/cgroups/subsystems"
	"github.com/niceforbear/docker-implementation-practice/container"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
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

func RunV2(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcessV2(tty)
	if parent == nil {
		logrus.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	// cgroup name
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()

	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(comArray, writePipe)

	parent.Wait()
}

func RunV3(tty bool, comArray []string) {
	parent, writePipe := container.NewParentProcessV2(tty)
	if parent == nil {
		logrus.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	sendInitCommand(comArray, writePipe)
	parent.Wait()
	os.Exit(0)
}

func RunV4(tty bool, comArray []string) {
	parent, writePipe := container.NewParentProcessV3(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	sendInitCommand(comArray, writePipe)
	parent.Wait()
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	container.DeleteWorkSpace(rootURL, mntURL)
	os.Exit(0)
}

func RunV5(tty bool, comArray []string, volume string) {
	parent, writePipe := container.NewParentProcessV4(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	sendInitCommand(comArray, writePipe)
	parent.Wait()
	mntURL := "/root/mnt"
	rootURL := "/root"
	container.DeleteWorkSpaceV2(rootURL, mntURL, volume)
	os.Exit(0)
}

// 初始化容器
func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}