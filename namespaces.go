package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// UTS NS 主要用来隔离 node name & domain name
func setupUtsNs() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// # pstree -pl
	// 查看进程关系

	// # echo $$
	// 输出当前PID

	// # readlink /proc/<PID>/ns/uts
	// 检查父进程 & 子进城是否不在同一个 UTS NS 下

	// # hostname -b <another_name>
	// # hostname
	// 修改 hostname，但是对宿主机 hostname 没影响
}

// 每一个 IPC NS 都有自己的 System V IPC & POSIX MQ
func setupIpcNs() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// # ipcs -q  // NS里执行查看IPCS
	// # ipcmk -Q // 外部创建IPC
	// # ipcs -q  // NS里再次查询
}

// 父进程为X的PID映射到NS中PID为1
// 此处 ps / top 查看会使用/proc中内容，因此还是会看到X的PID
func setupPidNs() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// # echo $$
}

// 隔离进程看到的挂载点视图
func setupMountNs() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// # mount -t proc proc /proc  // 把当前NS的内核proc mount到 /proc
	// # ls /proc
	// # ps -ef
}

// 隔离用户的用户组ID
func setupUserNs() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(-1)

	// # id
}

func setupNetNs() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(-1)

	// # ifconfig
	// 在宿主机 & 容器内 ifconfig 显示不同即网络处于隔离状态了
}
