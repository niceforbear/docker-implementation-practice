package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

/**
作用：初始化容器内容，挂载 proc 文件系统。运行用户指定程序。
返回：创建完成，容器开始运行。

此处已在容器内部执行。
这是本容器执行的第一个进程。
使用 mount 挂载 /proc 文件系统，才能查看当前进程资源。

MountFlags：

* MS_NOEXEC：在本文件系统中不允许运行其他程序
* MS_NOSUID：在本系统中运行程序的时候，不允许 set-user-ID or set-group-ID
* MS_NODEV： 自 Linux 2.4 以来，所有 mount 的系统都会默认设定的参数。

syscall.Exec：黑魔法！
调用 Kernel 的 int execve 系统函数。作用：执行当前 filename 对应的程序。
它会覆盖当前进程的镜像、数据、堆栈、PID等。
通过调用这个方法，将用户指定的进程运行起来，把最初的 init 进程替换掉。
因此当进入到容器内时，会发现容器内的第一个程序是我们指定的进程。
这也是目前 Docker 使用的容器引擎 runC 的实现方式之一。
*/
func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("cmd %v", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	argv := []string{command}

	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}

	return nil
}