package container

import (
	"os"
	"os/exec"
	"syscall"
)

/**
作用：创建 ns 隔离的容器进程。
返回：配置好隔离参数的 cmd 对象。

这里是父进程，即当前进程执行的内容。
1. /proc/self/exe 调用中，/proc/self/ 指的是当前运行进程自己的环境，exec 其实就是自己调用了自己，使用这种方式对创建出来的进程进行初始化。
2. args 是参数，init 是传递给本进程的第一个参数。本例子，会去调用 initCommand 去初始化进程的一些环境和资源。
3. 下面的 clone 参数就是 fork 一个新进程，并且使用了 ns 进行隔离。
4. 如果用户制定了 -ti, 需要修改标准输入输出。
*/
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}