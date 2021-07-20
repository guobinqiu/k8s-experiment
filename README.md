基于Kubernates发布一个go的web应用
---

Kubernates(又称k8s)是什么？抄概念没意思，通俗讲kubernates就是一个容器编排的工具。一个web应用会由许多服务组成，
比如数据库服务，缓存服务，或者是另外一个业务层面的服务，这么多服务在云原生架构的年代里就会忽悠你要把他们
都跑在各自的容器里，因为跑在容器里比直接跑在服务器上成本低，因为一台服务上可以装多个容器，为了高可用你也可以在
多个服务器上跑，本个服务器都跑上对方的副本，如何做到这些呢？这就是kubernates容器编排工具可以帮我们做到的，它
主要通过定义各种配置文件来调度容器内的服务。

要用kubernates发布一个go的web应用，主要经历三大过程

1. [搭建kubernates集群。集群的搭建有多种方式](./a.md)
   - 购买现成的各家云服务商已经整合好的集群服务
   - 在自己的电脑上搭建单机集群，推荐使用工具[minikube](https://minikube.sigs.k8s.io/docs/)+[virtualbox](https://www.virtualbox.org/)
   - 自己购买云服务器来搭建集群，推荐使用[kubeadmin](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)工具（本案例使用改方式）
   - 在本地机房搭建集群，推荐使用[kubeadmin](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)+[metallb](https://metallb.universe.tf/)
2. [容器化你本地的go应用](./b.md)，可以基于docker或者containerd，推荐containerd，containerd还没来得及学，本案例采用docker
3. [把go的web镜像跑在kubernates集群里](./c.md)

[服务进程目录](./d.md)
