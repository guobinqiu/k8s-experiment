# 搭建kubernates集群

## 生产环境

> https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/high-availability/

生产环境考虑到高可用性，一般至少需要9台服务器，3台跑`control-plane`，3台跑`worker`，3台跑`etcd`

- control-plane：相当于master节点
- worker：相当于slave节点，现在不能乱喊了
- etcd：持久化集群状态

## 测试环境

本案例考虑到个人学习成本问题，只从阿里云买了两台ec2云服务器，一台作为worker节点，另一台既作master又作worker

### 系统环境

系统|内网IP|公网IP|主机名|主机配置|节点类型
---|---|---|---|---|---
CentOS 7.9 64位|172.19.96.118|47.96.172.142|k8s-masternode|1核 2GiB|master & worker
CentOS 7.9 64位|172.19.44.93|47.98.221.22|k8s-worknode|1核 2GiB|worker

注意：官方推荐CPU至少2核，内存2G，由于我CPU是1核，在安装时需要忽略配置检查`kubeadm init ... --ignore-preflight-errors=NumCPU`

### 所有节点安装

#### 关闭防火墙

```
systemctl stop firewalld.service
systemctl disable firewalld.service
```
systemctl取代了早期的init.d

#### 禁用SELINUX

临时禁用
```
setenforce 0
```

永久禁用 
```
vim /etc/selinux/config
SELINUX=disabled
```

#### 修改k8s.conf文件

```
cat <<-EOF > /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF

sysctl --system
```
把EOF之间的内容输出到/etc/sysctl.d/k8s.conf文件

#### 关闭swap

临时关闭
```
swapoff -a
```
 
永久关闭

修改/etc/fstab文件，注释掉SWAP的自动挂载（永久关闭swap，重启后生效）

注释掉以下字段

```
/dev/mapper/cl-swap     swap                    swap    defaults        0 0
```

#### docker安装

```
yum install -y docker
```

#### [安装kubeadm、kubelet、kubectl](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/#installing-kubeadm-kubelet-and-kubectl)

- kubeadm: 安装k8s集群的工具
- kubelet: 管理着容器和pod的生命周期
- kubectl: k8s集群管理的客户端工具

```
sudo yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes
sudo systemctl enable --now kubelet
```

### Master节点安装

> https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/

#### 修改主机名

```
hostnamectl set-hostname k8s-masternode
```

#### 修改yum安装源

```
cat <<-EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
```
把EOF之间的内容输出到/etc/yum.repos.d/kubernetes.repo文件，使用阿里云源

#### 初始化一个只有主节点的k8s集群

```
kubeadm init \
--kubernetes-version=xxx
--apiserver-advertise-address=47.96.172.142 \
--image-repository registry.aliyuncs.com/google_containers \
--service-cidr=10.1.0.0/16 \
--pod-network-cidr=10.244.0.0/16
```
- kubernetes-version: 用于指定k8s版本
- apiserver-advertise-address: 用于指定kube-apiserver监听的ip地址,就是master本机IP地址
- pod-network-cidr: 用于指定Pod的网络范围：10.244.0.0/16
- service-cidr: 用于指定svc的网络范围
- image-repository: 指定阿里云镜像仓库地址，由于kubeadm默认从官网k8s.grc.io下载所需镜像，国内无法访问，因此需要通过–image-repository指定阿里云镜像仓库地址

集群初始化成功后返回如下信息：

```
Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

You should now deploy a Pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  /docs/concepts/cluster-administration/addons/

You can now join any number of machines by running the following on each node
as root:

  kubeadm join <control-plane-host>:<control-plane-port> --token <token> --discovery-token-ca-cert-hash sha256:<hash>
```

`kubeadm join <control-plane-host>:<control-plane-port> --token <token> --discovery-token-ca-cert-hash sha256:<hash>`
这行保存下来，将来要在worker节点上执行

#### 新增非root用户

为了能够在任何一台机器上都能管理k8s集群，需要在服务器上新增一个普通用户，并赋予sudo权限

```
adduser guobin
passwd guobin
usermod -aG wheel guobin
su - guobin
```

把guobin用户添加到wheel组就会自动拥有sudo权限，而不需要修改/etc/sudoers文件，这是一个小技巧

#### 配置kubectl

```
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

这样普通用户就可以在任意节点上使用kubectl命令了

#### [集群网络安装](https://kubernetes.io/docs/concepts/cluster-administration/networking/)

kubernates的CNI网络插件有很多，这里我们选择安装flannel，因为它的通用性比较好，如果你的k8s是搭建在自建机房的裸机上的话，有些网络插件会有不兼容的情况。如果你使用的是云服务商提供的k8s集群那他们一般都会有自己的CNI。

```
wget https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
```

如果yml中的"Network": "10.244.0.0/16"和kubeadm init xxx --pod-network-cidr不一样，就需要修改成一样的。不然可能会使得Node间Cluster IP不通。

由于我上面的kubeadm init xxx --pod-network-cidr就是10.244.0.0/16。所以此yaml文件就不需要更改了。

安装网络

```
kubectl apply -f kube-flannel.yml
```

### Worker节点安装

#### 修改主机名

```
hostnamectl set-hostname k8s-worknode
```

#### 新增非root用户

为了能够在任何一台机器上都能管理k8s集群，需要在服务器上新增一个普通用户，并赋予sudo权限

```
adduser guobin
passwd guobin
usermod -aG wheel guobin
su - guobin
```

把guobin用户添加到wheel组就会自动拥有sudo权限，而不需要修改/etc/sudoers文件，这是一个小技巧

#### 配置kubectl

```
mkdir -p $HOME/.kube
sudo scp guobin@47.96.172.142:~/.kube/config $HOME/.kube #把config文件从master节点复制过来
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

#### worker节点加入集群

登录到worker节点，确保已经安装了docker和kubeadm，kubelet，kubectl

执行你在master节点执行`kubeadm init ...`的时候保存下来的那串文本

```
kubeadm join <control-plane-host>:<control-plane-port> --token <token> --discovery-token-ca-cert-hash sha256:<hash>
```

### 最后

#### 登录任意节点查看Pod状态

```
kubectl get pod --all-namespaces -o wide
```

输出

```
NAMESPACE     NAME                                                       READY   STATUS    RESTARTS   AGE     IP              NODE           NOMINATED NODE   READINESS GATES
kube-system   coredns-7ff77c879f-gh8gf                                   1/1     Running   2          8d      10.244.0.7      k8s-master     <none>           <none>
kube-system   coredns-7ff77c879f-rjfkk                                   1/1     Running   2          8d      10.244.0.6      k8s-master     <none>           <none>
kube-system   etcd-k8s-master                                            1/1     Running   2          8d      172.19.96.118   k8s-master     <none>           <none>
kube-system   kube-apiserver-k8s-master                                  1/1     Running   2          8d      172.19.96.118   k8s-master     <none>           <none>
kube-system   kube-controller-manager-k8s-master                         1/1     Running   2          8d      172.19.96.118   k8s-master     <none>           <none>
kube-system   kube-flannel-ds-8xdvr                                      1/1     Running   1          8d      172.19.44.93    k8s-worknode   <none>           <none>
kube-system   kube-flannel-ds-tppxx                                      1/1     Running   2          8d      172.19.96.118   k8s-master     <none>           <none>
kube-system   kube-proxy-szzpx                                           1/1     Running   1          8d      172.19.44.93    k8s-worknode   <none>           <none>
kube-system   kube-proxy-td56t                                           1/1     Running   2          8d      172.19.96.118   k8s-master     <none>           <none>
kube-system   kube-scheduler-k8s-master                                  1/1     Running   2          8d      172.19.96.118   k8s-master     <none>           <none>
```

#### 登录到各节点查看IP

```
ifconfig
```

输出

```
flannel.1: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1450
        inet 10.244.0.0  netmask 255.255.255.255  broadcast 10.244.0.0
        inet6 fe80::78d7:82ff:fe8f:164  prefixlen 64  scopeid 0x20<link>
        ether 7a:d7:82:8f:01:64  txqueuelen 0  (Ethernet)
        RX packets 606351  bytes 36127936 (34.4 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 1574080  bytes 86788308 (82.7 MiB)
        TX errors 0  dropped 8 overruns 0  carrier 0  collisions 0
```

多出了一块flannel.1网卡，它是用来做集群内部网络通信的

这样kubernates集群搭建就ok了

## 补充

关于Ubuntu下的安装

1.更新软件源安装kubectl kubeadm kubelet containerd

```
curl -fsSL https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
 
sudo apt update

sudo apt install kubectl kubeadm kubelet containerd
```

2.生成一个配置集群的模板文件

```
kubeadm config print init-defaults > kubeadmin-config.yaml
```

修改后如下

```
apiVersion: kubeadm.k8s.io/v1beta3
bootstrapTokens:
- groups:
  - system:bootstrappers:kubeadm:default-node-token
  token: abcdef.0123456789abcdef
  ttl: 24h0m0s
  usages:
  - signing
  - authentication
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: 192.168.1.4
  bindPort: 6443
nodeRegistration:
  criSocket: unix:///var/run/containerd/containerd.sock
  imagePullPolicy: IfNotPresent
  name: node
  taints: null
---
apiServer:
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta3
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns: {}
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: registry.aliyuncs.com/google_containers
kind: ClusterConfiguration
#kubernetesVersion: 1.28.0
networking:
  dnsDomain: cluster.local
  serviceSubnet: 10.95.0.0/12
  podSubnet: 10.245.0.0/16
scheduler: {}
```
- advertiseAddress 我的虚拟机IP, apiserver连这个地址
- podSubnet Pod网络的地址范围,默认10.244.0.0/16, podSubnet需要与你的网络插件（如 Flannel）配置相匹配, 我有另外一个集群使用了默认配置, 为了不冲突改成10.245.0.0/16
- serviceSubnet Service的虚拟IP地址范围, 默认10.96.0.0/12, 我有另外一个集群使用了默认配置, 为了不冲突改成10.95.0.0/12
- imageRepository改成阿里云镜像registry.aliyuncs.com/google_containers

3.安装集群

```
sudo kubeadm init --config kubeadmin-config.yaml
```

安装失败

沙箱问题: 把`/etc/containerd/config.toml`文件里的`sandbox_image`替换成`registry.aliyuncs.com/google_containers/pause:3.9` (如果你当前是3.8, 它推荐你安装3.9, 虽然只是推荐但是不修改init不会成功)

kubelet启动失败: 可能启动参数不对, 从`/var/lib/kubelet/kubeadm-flags.env`文件删除`--container-runtime=remote`

worker节点加入集群 (可选, 单节点不需要执行)

```
kubeadm join 192.168.1.9:6443 --token abcdef.0123456789abcdef  --discovery-token-ca-cert-hash sha256:9f792067a16addee3a5f60150feb0289008db84c9d3711af3a4ce6fbcbd4f3a8
```

```
guobin@ubuntu02:~$ kubectl get nodes
NAME       STATUS   ROLES           AGE    VERSION
ubuntu02   Ready    control-plane   4d3h   v1.30.4
```

本地测试CNI用flannel网络

https://github.com/flannel-io/flannel/blob/master/README.md

```
net-conf.json: |
    {
      "Network": "10.245.0.0/16",
      "EnableNFTables": false,
      "Backend": {
        "Type": "vxlan"
      }
    }
```

- 修改Network, 它要和podSubnet保持一致
- image从docker.io改成docker.m.daocloud.io

最后

如果集群启动有问题, 尝试重启kubelet重新加载更新后的配置文件, 重启过程中各种错误先不用管, 等待一段时间再看

```
sudo systemctl restart kubelet
```
