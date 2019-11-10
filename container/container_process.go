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

func NewParentProcessV2(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// 额外携带的文件描述符。一般文件有 3 个 I、O、Error，通过查看 /proc/self/fd 可以看到多携带的文件描述符（也会以虚拟文件系统的方式看到）
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

/*
Pipe: Linux IPC 的一种实现
Pipe：一般是半双工，一端写，另一端读

管道类型：

* 无名管道：具有亲缘关系的进程之间使用
* 有名管道（FIFO管道）：存在于文件系统的管道。mkfifo() 创建

管道也是文件的一种，但是有固定 buffer，一般是 4KB。
管道满时，写进程会阻塞。读进程同理。
 */
func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}