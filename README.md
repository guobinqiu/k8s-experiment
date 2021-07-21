基于Kubernates发布一个go web应用
---

要用kubernates发布一个go web应用，主要经历3大过程

1. [搭建kubernates集群](install-cluster.md)有多种方式
   - 购买现成的各家云服务商已经整合好的集群服务
   - 在自己的电脑上搭建单机集群，推荐使用工具[minikube](https://minikube.sigs.k8s.io/docs/)+[virtualbox](https://www.virtualbox.org/)
   - 自己购买云服务器来搭建集群，推荐使用[kubeadmin](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)工具（本案例使用该方式）
   - 在本地机房搭建集群，推荐使用[kubeadmin](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)+[metallb](https://metallb.universe.tf/)
2. [把你本地的go web应用制成容器镜像](dockerize-go-app.md)
3. [把容器镜像跑在kubernates集群里](deploy-to-cluster.md)
