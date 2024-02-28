# Publish a go web application based on kubernates

[简体中文](README_CN.md)

To publish a go web application with kubernates, there are mainly three major processes

1. [Build kubernates cluster](install-cluster.md) There are many ways
    - Purchase ready-made cluster services that have been integrated by various cloud service providers
    - To build a stand-alone cluster on your own computer, it is recommended to use the tool [minikube](https://minikube.sigs.k8s.io/docs/)+[virtualbox](https://www.virtualbox.org/)(suitable for test only)
    - Purchase a cloud server by yourself to build a cluster. It is recommended to use the [kubeadmin](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/) tool (this method is used in this case)
    - To build a cluster in the local computer room, it is recommended to use [kubeadmin](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)+[metallb](https://metallb.universe.tf/)

2. [Make your local go web application into a container image](dockerize-go-app.md)
3. [Run the container image in the kubernates cluster](deploy-to-cluster.md)
4. [Install Kubernetes Dashboard](dashboard.md)

## License

MIT
