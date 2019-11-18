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

func RunV6(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcessV3(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	if tty {
		parent.Wait()
	}
}

// 初始化容器
func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func recordContainerInfo(containerPID int, commandArray []string, containerName string) (string, error) {
	id := randStringBytes(10)
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")
	if containerName == "" {
		containerName = id
	}
	containerInfo := &container.ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.RUNNING,
		Name:        containerName,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirUrl, err)
		return "", err
	}
	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}

	return containerName, nil
}

func RunV7(tty bool, comArray []string, res *subsystems.ResourceConfig, containerName string) {
	parent, writePipe := container.NewParentProcessV3(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	//record container info
	containerName, err := recordContainerInfo(parent.Process.Pid, comArray, containerName)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	if tty {
		parent.Wait()
		deleteContainerInfo(containerName)
	}
}

func deleteContainerInfo(containerId string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s error %v", dirURL, err)
	}
}