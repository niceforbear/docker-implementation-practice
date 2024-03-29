# 进阶

## 容器的后台运行

* 实现：`docker -d`
* 最开始，所有容器 init 进程都是从 docker daemon 进程 fork 出来。
* 会出现一个问题：daemon 挂掉，所有容器挂掉。
* 后来，Docker 使用了 containerd，即现在的 runc，可以实现 daemon 挂掉，但是容器依然存在。

* 因此需要父进程创建完成子进程后，`detach` 掉子进程。

实现原理

* 容器是 docker fork 出来的子进程。
* 如果父进程推出，子进程称为`孤儿进程`。
* 为了避免孤儿进程退出时无法释放所占用的资源而僵死，进程号为1的进程 init 会接受孤儿进程。

* 如果 detach 创建了容器，就不能再等待。因此，将容器内 init 进程启动起来后可以退出，然后由 OS 进程 ID 为 1 的 init 进程接管容器进程。

## 实现容器查看

* 实现：`docker ps`
* 需要做到：指定容器 name
* 记录容器信息到文件，存储到宿主机上

## 实现查看容器日志

* 实现: `docker logs`
* 修改：在创建容器时，调整标准输出重定向到 log 文件
* 执行 `docker logs` 时，获取 log 文件内容。 

## 实现进入容器 Namespace

* 实现：`docker exec`
* `setns`: 通过这个系统调用，根据提供的 PID 进入指定的 Namespace
    * 原理：打开 /proc/<pid>/ns 下的对应文件
* 问题：对于 Mount Namespace，一个具有`多线程的进程`是无法使用 setns 调用进入到对应的 NS 的。
* 因为 Go `每启动一个程序就会进入多线程状态`，因此无法简单通过 Go 直接进行系统调用
* 需要借助 Cgo
* `nsenter.go`
    * 作用：一旦包被导入，则会在所有 Go 运行的环境之前执行，避免 Go 多线程无法进入 mnt NS
    * 问题：会在所有 Go 代码前执行，即使不需要使用 exec 的代码。e.g. `mydocker run`
        * 解决：在开始位置指定环境变量，在需要使用的地方注入环境变量来使用。

* 当通过容器名和 cmd 进来后，Go 程序已经执行，并且已经运行 C 代码（因为会 load package 进行加载）。
* 通过使用 /proc/self/exec 来再执行一遍。
* 只需要简单 fork 一个进程，不需要 NS 和隔离，将这个 fork 的进程的 标准输入/输出 bind 到宿主机。
* 由于第二次运行制定了环境变量，因此可以进入指定的 NS。 

## 实现停止容器

* stop 容器的原理：查找它的主进程 PID，发送 SIGTERM 信号，等待进程结束。

## 容器删除

* 实现 `rm` 操作。
* `rm` 操作在已经关闭 container 情况下，只需要把记录文件信息的内容删除即可。

## 通过容器构建镜像

* 使用 `docker run -d` 能够同时运行多个容器，这样会存在一个问题：使用一个镜像启动多个容器会共用一个 AUFS 文件系统。

* 此节要实现：
    * 每个容器单独分配隔离的文件系统
    * 修改 commit 命令，实现对不同容器可以打包镜像

## 实现容器指定环境变量

* 新启动的进程默认继承父进程的环境变量
* 将自定义的环境变量 append 到 os.Environ 即可。
* exec 执行后看到的变量实际上是宿主机的环境变量
* 在使用 exec 时，通过 /proc/PID/environ 获取 NS 的环境变量，生成到 exec 的 NS 中。

## 