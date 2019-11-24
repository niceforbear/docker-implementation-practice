# Network

## 容器网络虚拟化

* Linux 实际是通过网络设备去操作和使用网卡的，系统装了一个网卡后，会为其生成一个网络设备实例，e.g. `eth0`。
* Linux 常见的虚拟化设备有：Veth, Bridge, 802.1q VLAN device, TAP.

### Linux Veth

* Veth 是成对出现的虚拟网络设备，发送到 Veth 一端虚拟设备的请求会从另一端的虚拟设备中发出。
* 在容器场景中，使用 Veth 连接不同的网络 NS。

```bash
# 创建两个网络 NS
$ sudo ip netns add ns1
$ sudo ip netns add ns2

# 创建一对 Veth
$ sudo ip link add veth0 type veth peer name veth1

# 将两个 Veth 移到两个 NS 中
$ sudo ip link set veth0 netns ns1
$ sudo ip link set veth1 netns ns2

# 在 ns1 的 NS 中查看网络设备
$ sudo ip netns exec ns1 ip link
# 此时会看到除了 loopback 之外的网络设备 veth

# 给两端分配不同的地址后，向 veth 一端发送请你去，就能达到另一端

# 配置每个 veth 的网络地址和 NS 路由
$ sudo ip netns exec ns1 ifconfig veth0 172.18.0.2/24 up
$ sudo ip netns exec ns2 ifconfig veth1 172.18.0.3/24 up
$ sudo ip netns exec ns1 route add default dev veth0
$ sudo ip netns exec ns2 route add default dev veth1

# 发包测试
$ sudo ip netns exec ns1 ping -c 1 172.18.0.3
``` 

### Linux Bridge

* Bridge 虚拟设备是用来桥接的网络设备，相当于现实世界中的交换机，可以连接不同的网络设备。
* 当请求到达 Bridge 时，可以通过报文中的 Mac 地址进行广播或转发。
* 例子如下：创建一个 Bridge，连接 NS 中的网络设备和宿主机的网络。

```bash
# 创建 veth 设备并将一端移入 NS
$ sudo ip netns add ns1
$ sudo ip link add veth0 type veth peer name veth1
$ sudo ip link veth1 setns ns1

# 创建 Bridge，挂载网络设备
$ sudo brctl addbr br0
$ sudo brctl addif br0 eth0
$ sudo brctl addif br0 veth0
```

eth0 <-> br0 <-> veth0 <-> veth1

### Linux 路由表

* 路由表是 Linux kernel 的一个模块，通过定义路由表来决定在某个网络 NS 中包的流向，从而定义请求会到哪个网络设备上。
```bash
# 启动虚拟网络设备，并设置它在 Net NS 中的 IP 地址。
$ sudo ip link set veth0 up
$ sudo ip link set br0 up
$ sudo ip netns exec ns1 ifconfig veth1 172.18.0.2/24 up

# 分别设置 ns1 的路由和宿主机上的路由
# default 代表 0.0.0.0/0, 即在 Net NS 中的所有流量都经过 veth1 的网络设备流出
$ sudo ip netns exec ns1 route add default dev neth1

# 在宿主机上将 172.18.0.0/24 的网段请求路由到 br0 的网桥
$ sudo route add -net 172.18.0.024 dev br0

# 查看宿主机的 IP 地址查询
$ sudo ifconfig eth0

# 从 NS 中访问宿主机的地址
$ sudo ip netns exec ns1 ping -c x.x.x.x

# 从宿主机访问 NS 中的地址
$ ping -c 1 172.18.0.2
```

### Linux iptables

* `iptables` 是对 Linux 内核的 netfilter 模块进行操作和展示的工具，管理包的流动和转送。
* `iptables` 定义了一套链式处理的结构，在网络报传输的各个阶段可以使用不同的策略对包进行加工、传送或丢弃。
* 在容器网络虚拟化中，常用两种策略用于容器和宿主机外部的网络通信：`MASQUERADE` & `DNAT`

#### MASQUERADE

* 该策略可以将请求包中的源地址转换成另一个网络设备的地址。
* e.g. NS 中网络设备的 IP 是 172。18。0。2，在宿主机上可以路由到 br0 网桥，但是到达宿主机外部后，是不知道如何路由到这个 IP 的，因此请求外部地址，需要先通过此策略将 IP 转换成宿主机出口网卡的 IP。
```bash
# 打开 IP 转发
$ sudo sysctl -w net.ipv4.conf.all.forwarding=1
# 对 NS 中发出的包添加网络地址转换
$ sudo iptables -t nat -A POSTROUTING -s 172.18.0.0/24 -o eth0 -j MASQUERADE
```

#### DNAT

* 此策略也是做网络地址转换的，不过其目的是更换目标地址，经常用于将内部网络地址的端口映射到外部去。
* e.g. NS 如果需要提供服务给宿主机之外的应用。外部没法使用 172 的这个 IP，需要用到 DAT。

```bash
# 将宿主机上 80 端口的 TCP 请求转发到 NS 的 IP:80 上. 
$ sudo iptables -t nat -A PREROUTING -p tcp -m tcp --dport 80 -j DNAT --to-destination 172.18.0.2:80
```

### Go Network library

* `net` library
* github.com 上开源库, e.g. netns library

## 容器网络模型

* 容器网络的两个对象：网络 & 网络端点
    * 网络 network：容器的一个集合，在这个网络上的容器可以通过这个网络互相通信，参考 Bridge。
    * 网络端点 endpoint：
        * 连接容器与网络的，保证容器内部与网络的通信。参考 Veth。
        * 网络端点的传输依靠：
            * 网络驱动 Network driver：网络功能中的组件。不同的驱动对网络的创建、连接、销毁的策略不同，通过创建网络时指定不同的网络驱动来定义使用哪个配置
            * IPAM：网络功能的组件。用于网络 IP 地址的分配和释放，包括容器的 IP 地址和网络网关的 IP，主要功能包括：
                * IPAM.Allocate(subnet *net.IPNet)：从指定的 subnet 网段中分配 IP 地址。
                * IPAM.Release(subnet net.IPNet, ipaddr net.IP)：从指定的 subnet 网段中释放掉指定的 IP 地址。

### 调用关系

#### 创建网络

```bash
mydocker network create --subnet 192.168.0.0/24 --driver bridge testbridgenet
```

#### Actions

```bash
# connect
mydocker runn -it -p 80:80 --net testbridgenet xxxx

# list
mydocker network list

# remove
mydocker network remove testbridgennet

```

## 容器地址分配

### bitmap 算法

* 在大规模连续且少状态的数据处理中有很高的效率，e.g. IP 地址分配。使用位图进行管理。
* 在网段中，某个 IP 地址有两个状态：1 - 被分配，0 - 未分配。

## Bridge 网络

* 在网络挂载 endpoint 实现容器的互相通信以及容器外部通信

## 跨主机容器网络

### 跨主机容器网络的 IPAM

* 如果每个机器只负责容器网络在自己宿主机上的 IP 地址分配，那么可能会造成不同机器上分配到的容器 IP 地址重复的问题，进而导致访问问题。
* 通常采用中心化的一致性 KV-store 来记录 IP 地址的分配存储。
* 当多个宿主机同时访问时会引起并发问题，通常使用以下两种方法时下节点一致：
    * 全局锁
    * Compare And Swap（CAS）
        * 通过 CAS 方法，每个宿主机在写入 IP 分配的信息时，会判断分配过程中这个 IP 分配信息的 key 是否修改过。
        * 如果修改过，则重新获取信息并重新分配 IP 后重试
        
### 跨主机容器网络常见实现方式

* 封包
    * 宿主机之前不知道容器的地址怎么路由和访问，可以把容器之间的请求外面包装上宿主机的地址发送
    * 常见封包技术：Vxlan，Ipip-tunnel，GRE
* 路由
    * 让宿主机的网络知道容器的地址怎么路由，以及路由到哪个机器
    * 这种方式一般需要网络设备的支持，e.g. 修改路由器的路由表，将容器 IP 地址的下一跳修改到这个容器所在的宿主机，来达到跨主机容器请求
    * 常见技术：路由器路由表，Vlan，VPC 路由表
    
对比

* 封包
    * 基础设施要求低
    * 性能损耗大
* 路由
    * 性能好
    * 基础设施有要求，需要支持路由的配置或者特定网络协议