package subsystems

/*
* cgroup hierarchy 中的节点：管理进程和 subsystem 的控制关系。
* subsystem：作用于 hierarchy 上的 cgroup 节点，控制节点中进程的资源占用。
* hierarchy：将 cgroup 通过树状结构串起来，通过虚拟文件系统方式暴露。
*/

type ResourceConfig struct {
	MemoryLimit string
	CpuShare	string
	CpuSet		string
}

/**
Subsystem 接口，每个 subsystem 可以实现下面的 4 个接口。
此处将 cgroup 抽象为 path：因为 cgroup 在 hierarchy 的路径，就是虚拟文件系统中的虚拟路径
*/
type Subsystem interface {
	// subsystem 名字，e.g. cpu memory
	Name() string
	Set(path string, res *ResourceConfig) error

	// 将进程添加到某个 cgroup 中
	Apply(path string, pid int) error
	Remove(path string) error
}

var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
