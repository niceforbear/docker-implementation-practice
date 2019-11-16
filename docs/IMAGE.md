# Build Image

## 通过 busybox 创建容器镜像基础

Image 使得容器传递和迁移更加简单。

使用 busybox 文件系统的 rootfs，使用 `docker export` 将镜像打成 tar 包。

```bash
docker pull busybox
docker run -d busybox top -b
docker export -o busybox.tar <container_id>
tar -xvf busybox.tar -C busybox/
```

### pivot_root

使用 pivot_root 改变 root 的文件系统。

* 系统调用，作用是改变当前 root 文件系统。
* pivot_root 可以使整个系统都切换到新的 root 目录。能够 unmount 之前的 root 文件系统。
* chroot 只能修改单个进程的 root 文件系统

### tmpfs

* tmpfs 是基于内存的文件系统，可以使用 RAM ／ swap 分区来存储。

## 使用 AUFS 包装 busybox 实现容器和镜像分离

* Docker 在使用镜像启动时，会新建两个 layer：
    * write layer
        * 容器唯一的读写层
    * container-init layer
        * 为容器新建的只读层
        * 存储容器启动时传入的系统信息
* 最后把 write layer & container-init layer 和相关镜像都 mount 到一个 mnt 目录下
* 把 mnt 目录作为容器启动的根目录

## 通过 volume 数据卷，实现容器内数据持久化

* 对 volume 进行 mount ／ unmount 操作。

## Build 镜像原理

* 实现 commit command
* 将多个 layer 目录一起打成 tar 包

