### 把容器镜像跑在kubernates集群里

先看一下发布所使用的配置文件的目录结构

```
Guobins-MBP:k8s-deployments guobin$ tree
.
├── app
│   ├── app-deployment.yml
│   └── app-service.yml
├── baremetal
│   ├── metallb.yml
│   └── values.yaml
├── haproxy
│   ├── haproxy-ingress-values.yaml
│   └── ingress.yml
├── mysql
│   ├── mysql-deployment.yml
│   ├── mysql-pv.yml
│   ├── mysql-pvc.yml
│   ├── mysql-secret.yml
│   ├── mysql-service.yml
│   └── mysql-service.yml.bak
└── net
    └── kube-flannel.yml

5 directories, 13 files
```

- app: go web服务的配置文件
- baremetal: 我实验用的，我们这里不会用到，这在本地机房部署k8s集群才会用
- haproxy: 灰度发布或者想省负载均衡器时用，我们这里用了，但只占个位，做个简单的转发，方便将来扩展
- mysql:我们的web服务依赖的数据库服务
- net: 集群网络配置

kubernates有几种对象：deployment、service、ingress、pv、pvc、secret等等，对应到我们这里的配置文件

对象：
- deployment：创建一组[pods](https://kubernetes.io/docs/concepts/workloads/pods/#what-is-a-pod)
- [service](https://kubernetes.io/docs/concepts/services-networking/service/)：暴露一组pods供外部访问
  - ClusterIP：internet不能访问，只能进到安装k8s的宿主机里去访问
  - NodePort：internet能访问，通过<NodeIP>:<NodePort>访问，本案例使用的此种方式
  - LoadBalancer: internet能访问的前提是启用云服务商提供的负载均衡服务
  - ExtenalName：没用过
- [ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)：暴露一组services供外部访问
- [pv](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)：数据持久化卷（声明）
- [pvc](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)：数据持久化卷（使用）
- [secret](https://kubernetes.io/docs/concepts/configuration/secret/)：类似环境变量

以下所有的发布都到默认的namespace：`default`，namespace可以用来区分发布环境，如`namespace=staging, namespace=dev`等

#### 部署go web服务

###### 创建一个go web deployment

从qiuguobin的dockerhub上拉取镜像部署到集群

```
kubectl apply -f app-deployment.yml
```

###### [创建一个go web service](https://kubernetes.io/docs/concepts/services-networking/service/)

暴露deployment创建的pods供集群外部访问，暴露端口为：30000（NodePort方式）

```
kubectl apply -f app-service.yml
```

#### 部署mysql服务

###### 创建一个mysql本地磁盘卷

```
kubectl apply -f mysql-pv.yml
kubectl apply -f mysql-pvc.yml
```

###### 创建一个mysql数据库环境变量文件

```
kubectl apply -f mysql-secret.yml
```

###### 创建一个mysql deployment

```
kubectl apply -f mysql-deployment.yml
```

###### [创建一个mysql service](https://kubernetes.io/docs/concepts/services-networking/service/)

不允许外部访问，只能集群内部之间访问（ClusterIP方式）

```
kubectl apply -f mysql-service.yml
```

###### mysql的彻底删除

除了删除持久化卷

```
kubectl delete pvc mysql-pv-claim
kubectl delete pv mysql-pv-volume
```

还要进到你的服务器里去删除其挂载目录

```
rm -rf /mnt/data
```

#### [部署haproxy ingress服务](https://github.com/haproxy-ingress/charts)

haproxy的配置文件有两种创建方式：deployment和[helm](https://helm.sh/docs/intro/install/)，我们这里采用helm方式来创建，这比自己写deployment省心很多哦

通过下面的脚本安装应该是最方便的

```
$ curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
$ chmod 700 get_helm.sh
$ ./get_helm.sh
```

我在阿里云服务器上下载不下来，所以是通过手动下载到本地，再上传到服务器上去编译安装的 
>https://helm.sh/docs/intro/install/#from-the-binary-releases

####### 更新helm源

```
helm repo add stable https://charts.helm.sh/stable
helm repo add incubator https://charts.helm.sh/incubator
```
stable库里是已经稳定的，incubator库里是正在孵化的

验证一下是否已加入

```
helm repo list

NAME     	URL
stable   	https://charts.helm.sh/stable
incubator	https://charts.helm.sh/incubator
```

###### 使用helm安装haproxy ingress服务

```
helm install haproxy incubator/haproxy-ingress --create-namespace --namespace default -f haproxy-ingress-values.yaml
```
haproxy-ingress-values.yaml这个文件告诉helm，你希望它怎么安装你的haproxy ingress服务。
这里我们让haproxy暴露一个30001端口供外部internet访问

###### 创建一个haproxy的转发规则

```
kubectl apply -f ingress.yml
```

这里直接转发所有外部流量到后端服务，未做任何分流处理，我们希望

#### 优化

生产环境我们访问网站都会使用80和433端口，这样可以省略口端号，而k8s的NodePort给我们暴露
的端口规定要从30000开始，为了隐藏NodePort暴露的端口，我申请了一个负载均衡SLB服务，IP为`121.40.221.176`，
这样所有对80端口的访问都可以被映射到30001端口。

在浏览器中当访问`http://121.40.221.176/`时，会输出
```
My Awesome Go App!!!
```

当访问`http://121.40.221.176/users`时，会输出
```
name=guobin, age=18
name=jack, age=100
```
这个数据是在mysql容器里手动写入的，由于是实验环境，我的mysql配置成了只能在集群内部访问，
你也可以把mysql服务配置成NodePort方式，这样你就可以在外面直接连接到容器内部的mysql去操作了。

