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

## 容器网络模型

## 容器地址分配

## Bridge 网络

## 跨主机容器网络
