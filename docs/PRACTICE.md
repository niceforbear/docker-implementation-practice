# Practice 实践

## runC 容器运行引擎

* 过去5年，Linux 逐步增加了 Cgroups，Namespace，Seccomp，capability，Apparmor 等功能，这些特性使得容器技术发展。
* Docker 公司将这些底层技术合并在一起，开源了 runC

runC 基本功能

* 完全支持 Linux NS
* 原生支持所有 Linux 的安全特性
* 在 CRIU 项目的支持下原生支持容器热迁移
* 一份正式的容器标准。

### OCI 标准包，bundle

* 一个标准的容器运行时需要文件系统，即`镜像`
* 组成
    * config.json：容器的配置。这个问价那必须在容器的 root 文件系统内。
    * 一个文件夹，代表容器的 root 文件系统。一般命名为 rootfs。这个文件夹必须包含上述 config.json

#### config.json

* ociVersion: OCI 容器的版本号
* root：配置容器的 root 文件系统
    * path：指定 root 文件系统的路径
    * readonly：为 true 的话，那么 root 文件系统在容器内是只读的，默认 false。

#### others

* mounts
    * 配置额外的挂载点
* process
    * 配置容器进程信息
* user
    * 指定容器内运行进程的用户信息
* hostname
* platform
    * os
    * arch
* hook
    * 可以扩展容器运行时的动作
    * 可以在容器运行前、停止后执行一些动作，e.g. 网络配置、垃圾回收
    
## 创建容器流程

* `runc run container-id` 根据 config.json 创建一个容器
* 读取配置文件
* 设置 rootFileSystem
* 使用 factory 创建容器，各个 OS 实现不同
* 创建容器的初始化进程 process
* 设置容器的输出管道，主要是 Go 的 pipes
* 执行 Container.Start() 启动物理的容器
* 回调 init 方法重新初始化容器进程
* runC 父进程等待子进程初始化成功后退出

## Docker containerd

* `containerd` 可以作为 daemon 程序运行，管理机器上所有容器的生命周期。
* containerd 主要集成于上层系统，e.g. Swarm，K8s
* 其通过 unix domain socket 暴露底层的 gRPC API，上层系统通过这些 API 管理机器上的容器。
* 每个 containerd 只负责一台机器，包括 pull image，对容器操作（start，stop），网络，存储
* 具体运行容器由 runC 负责执行。

### 架构

containerd 支持的能力如下：

* Distribution：和 Docker registry 打交道
* Bundle：管理本地磁盘上镜像的子系统
* Runtime：创建容器，管理容器的子系统

```text
             gRPC API | Metrics API
-----------------------------------------------------------
Subsystems | Distribution Bundle Runtime
-----------------------------------------------------------
Components | Content Metadata Snapshots Executor Supervisor
```

### 关系

* containerd 只负责 runtime 管理，而 Docker 还包括镜像构建。
* containerd 提供 偏底层API

* OCI：标准化的容器运行规范，包括 runtime 规范 & 镜像规范
* runC 是 OCI 的一个实现，Docker 贡献了 runC 的主要代码

* K8s 目前使用 Docker，未来可以使用 containerd

## K8s CRI 容器引擎

* CRI：容器运行时接口，是一组接口规范
* 此规范使得 K8s 不需要重新编译就可以使用更多的 runtime
* CRI 包括：PB，gRPC API 及 runtime library 支持。

CRI 概览

* kubelet 通过 gRPC 与 CRI shim 通信，CRI shim 通过 Unix Socket 启动一个 gRPC server 提供 runtime 服务
* kubelet 作为 gRPC client，通过 Socket 与 CRI shim通信。
* gRPC server 使用 pb 提供两类 gRPC service：ImageService & RuntimeService
* ImageService 提供镜像服务
* RuntimeService 提供容器相关生命周期管理及交互

接口与实现

* CRI 核心概念：PodSandbox * container
* Pod 由一组容器组成，这些容器共享环境与资源约束，这个共同的环境与资源约束称为 `PodSandbox`。
* 不同的 runtime 对 PodSandbox 的实现不同，因此 CRI 留有一组接口给不同的 runtime自主发挥，e.g. Hypervisor 将 PodSandbox 实现成虚拟机，Docker 将 PodSandbox 实现成一个 Linux NS。

Service 简介

* RuntimeService
    * Kubelet 在创建一个 Pod 前首先调用 RuntimeService
    * RunPodSandbox 为 Pod 创建一个 PodSandbox，这个过程包括：初始化 Pod 网络，分配 IP，激活 sandbox 等。
    * 然后 Kubelet 调用 CreateContainer，StartContainer，StopContainer，RemoveContainer 对容器进行操作
    * 当 Pod 删除时，会调用 StopPodSandbox，RemovePodSandbox 来销毁 Pod
    * Kubelet 职责在于 Pod 的生命周期的管理，包含健康检测、重启策略控制，同时实现你容器生命周期中的各种 hook
    
    * RuntimeService 还定义了与 Pod 中容器进行交互的接口
    * 目前 Kubelet 使用容器本地方法或 `nsenter/socat` 两种方法与 NS 交互。
* ImageService
    * CRI 中 interface
* LogService
    * 定义了容器的 stdout/stderr 应该如何被 CRI 处理的规范
    * CRI 通过以下 2 方面解决日志处理的问题
        * 强制指定 log 存放在本地文件的位置，并且 kubelet 和 日志收集器能方便地直接访问，e.g. `/var/log/pods/<podUID>/<containerName>_<instance#>.log`
        * runtime 需要按照 kubelet 能够理解的方式输出 log。
        
Why CRI？

* 不是所有的 runtime 都原生支持 Pod。为了支持所有 Pod 特性，当这些 runtime 与 K8s 集成时，需要大量时间来实现一个 shim
* 高层次的 interface 让代码共享、重用变得困难
* Pod Spec 演化速度很快。
