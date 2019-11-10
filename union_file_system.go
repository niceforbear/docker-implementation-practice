package main

// Union File System, i.e. UnionFS: 把其他文件系统联合到一个联合挂载点的文件系统服务
// 使用 branch 把不同文件系统的文件/目录进行覆盖，形成一个单一一致的文件系统。
// 这些 branch 是： read-only / read-write
// 看起来 UnionFS 可以对任何文件进行操作，但是并没有改变原来的文件，因为用到了 CoW

// CoW i.e. copy-on-write 写时复制：
// 如果一个资源是重复的，但没有任何修改，这是不需要立即创建一个新的资源，这个资源可以被新旧实例共享。
// 创建新资源在第一次写操作，即对资源进行修改时。

// AUFS i.g. Advanced Multi-Layered Unification Filesystem
// 为了可靠性/性能，重写了早期 UnionFS 1.x。同时引入新功能：e.g. 可写分支的 LB。

// AUFS 是 Docker 选用的第一种存储驱动。下面介绍 Docker 如何利用 AUFS 存储 image & container：

// image layer & AUFS
// 每一个 image 都是由一系列 read-only layer 组成。
// image layer 的内容都存储在 Docker hosts filesystem 的 /var/lib/docker/aufs/diff下
// /var/lib/docker/aufs/layers 存储着 image layer 如何堆栈这些 layer 的 metadata
// # docker history <image/tag>  // 查看image layer

// container layer & AUFS
// CoW 意味着：一旦某个文件只有很小的部分有改动，AUFS 也需要复制整个文件，这种设计对容器性能有影响。
// 但对于容器，每个 image layer 最多只需要复制一次，后续的改动会在第一次拷贝的 container layer 上进行。
// 启动一个 container 时，Docker 会为其创建一个 read-only 的 init layer，存储与这个容器内环境有关的内容；
// 同时为其创建一个 read-write layer来执行所有写操作。
// container read-write layer 会进行存储，因此容器停止/重启时，layer 依然存在。只有容器被删除时，read-write layer 才会一起删除。

// 可以尝试使用 AUFS & CoW 实现文件管理

