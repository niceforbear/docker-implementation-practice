package cgroups

import (
	"github.com/Sirupsen/logrus"
	"github.com/niceforbear/docker-implementation-practice/cgroups/subsystems"
)

// 将不同的 subsystem 中的 cgroup 管理起来，与容器建立关系。
// 作用：通过 CgroupManager，将资源限制的配置，以及进程移动到 cgroup 中的操作，通过各个 subsystem 去处理。

type CgroupManager struct {
	// cgroup 在 hierarchy 中的路径，相当于创建的 cgroup 目录相对于 root cgroup 目录的路径
	Path     string
	// 资源配置
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// 将进程 pid 加入到这个 cgroup 中
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// 设置 cgroup 资源限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

//释放 cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
