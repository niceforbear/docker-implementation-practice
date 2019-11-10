package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

// Linux Cgroups: 对一组进程及将来子进程的：资源限制 / 控制 / 统计
// 资源包括：CPU，内存，存储，网络

// Cgroups 三组件：
// * cgroup：对进程分组管理的一种机制。一个 cgroup 包含一组进程，并可以在这个 cgroup 上增加 Linux subsystem 的各种参数配置，将一组进程和一组 subsystem 的系统参数关联。
// * subsystem：一组资源控制的模块。
// * hierarchy：hierarchy 将 cgroup 串成树状结构，这样，Cgroups 就可以做到继承。

// 组件关系：
// 系统创建了新的 hierarchy 后，所有进程都会加入这个 hierarchy 的 cgroup 根结点
// 一个 subsystem 只能附加到一个 hierarchy 上面
// 一个 hierarchy 可以附加多个 subsystem
// 一个进程可以作为多个 cgroup 的成员，但是这些 cgroup 必须在不同的 hierarchy 中
// 一个进程 fork 出子进程时，子进程 & 父进程在同一个 cgroup 中

// Kernel 接口

// 1. 创建并挂载一个 hierarchy， i.e. cgroup 树
// 2. 在刚创建好的 hierarchy 上 cgroup 根结点中扩展两个子 cgroup
// 3. 在 cgroup 中添加 / 移动进程
// 4. 通过 subsystem 限制 cgroup 中进程的资源

const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func setupCgroupMemory() {
	if os.Args[0] == "/proc/self/exe" {
		// 容器进程
		fmt.Printf("current pid %d", syscall.Getpid())
		fmt.Println()

		cmd := exec.Command("sh", "-c",`stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	} else {
		fmt.Printf("%v", cmd.Process.Pid)

		// 在系统默认创建挂载了 memory subsystem 的 Hierarchy 上创建 cgroup
		_ = os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit"), 0755)

		// 将容器进程加入到这个 cgroup 中
		_ = ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)

		// 限制 cgroup 进程使用
		_ = ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)

		cmd.Process.Wait()
	}
}
