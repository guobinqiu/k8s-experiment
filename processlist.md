### 服务进程目录

#### All Pods

```
[guobin@k8s-master ~]$ kubectl get pods --all-namespaces
NAMESPACE     NAME                                                       READY   STATUS    RESTARTS   AGE
default       go-web-app-deployment-6fd8d76dff-gdmbn                     1/1     Running   0          4d2h
default       go-web-app-deployment-6fd8d76dff-ww6g5                     1/1     Running   0          4d2h
default       haproxy-haproxy-ingress-controller-697c5bc66c-cnk6d        1/1     Running   0          24h
default       haproxy-haproxy-ingress-controller-697c5bc66c-wlp4j        1/1     Running   0          24h
default       haproxy-haproxy-ingress-default-backend-5b74fff5f7-gtxrs   1/1     Running   0          24h
default       mysql-6db984b79d-jq7qq                                     1/1     Running   0          4d23h
kube-system   coredns-7ff77c879f-gh8gf                                   1/1     Running   2          8d
kube-system   coredns-7ff77c879f-rjfkk                                   1/1     Running   2          8d
kube-system   etcd-k8s-master                                            1/1     Running   2          8d
kube-system   kube-apiserver-k8s-master                                  1/1     Running   2          8d
kube-system   kube-controller-manager-k8s-master                         1/1     Running   2          8d
kube-system   kube-flannel-ds-8xdvr                                      1/1     Running   1          8d
kube-system   kube-flannel-ds-tppxx                                      1/1     Running   2          8d
kube-system   kube-proxy-szzpx                                           1/1     Running   1          8d
kube-system   kube-proxy-td56t                                           1/1     Running   2          8d
kube-system   kube-scheduler-k8s-master                                  1/1     Running   2          8d
```

#### All Deployments

```
[guobin@k8s-master ~]$ kubectl get deployments
NAME                                      READY   UP-TO-DATE   AVAILABLE   AGE
go-web-app-deployment                     2/2     2            2           4d2h
haproxy-haproxy-ingress-controller        2/2     2            2           24h
haproxy-haproxy-ingress-default-backend   1/1     1            1           24h
mysql                                     1/1     1            1           4d23h
```

#### All Services

```
[guobin@k8s-master ~]$ kubectl get services
NAME                                      TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)                      AGE
go-web-service                            NodePort    10.1.12.74    <none>        8080:30000/TCP               4d
haproxy-haproxy-ingress-controller        NodePort    10.1.48.26    <none>        80:30001/TCP,443:31539/TCP   24h
haproxy-haproxy-ingress-default-backend   ClusterIP   10.1.25.132   <none>        8080/TCP                     24h
kubernetes                                ClusterIP   10.1.0.1      <none>        443/TCP                      8d
mysql                                     ClusterIP   None          <none>        3306/TCP                     4d2h
```

#### All Ingresses

```
[guobin@k8s-master ~]$ kubectl get ingress
NAME                     CLASS    HOSTS   ADDRESS   PORTS   AGE
go-web-service-ingress   <none>   *                 80      24h
```

#### All PVs

```
[guobin@k8s-master ~]$ kubectl get pv
NAME              CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                    STORAGECLASS   REASON   AGE
mysql-pv-volume   1Gi        RWO            Retain           Bound    default/mysql-pv-claim   manual                  4d23h
```

#### All PVCs

```
[guobin@k8s-master ~]$ kubectl get pvc
NAME             STATUS   VOLUME            CAPACITY   ACCESS MODES   STORAGECLASS   AGE
mysql-pv-claim   Bound    mysql-pv-volume   1Gi        RWO            manual         4d23h
```

#### All Secrets

```
[guobin@k8s-master ~]$ kubectl get secret
NAME                                  TYPE                                  DATA   AGE
default-token-89dp6                   kubernetes.io/service-account-token   3      9d
haproxy-haproxy-ingress-token-vvtrm   kubernetes.io/service-account-token   3      46h
mysql-secret                          Opaque                                4      5d22h
sh.helm.release.v1.haproxy.v1         helm.sh/release.v1                    1      46h
```
