#基础镜像
FROM golang

#由于使用的不是go mod，GO111MODULE要关闭
ENV GO111MODULE=off

#国内下载
ENV GOPROXY=https://goproxy.cn

#下载govendor包
RUN go get -u -v github.com/kardianos/govendor

#在容器内创建/go/src/go-app目录，并把宿主机当前目录下所有文件拷贝至其中
ADD . /go/src/go-app

#容器内cd到上一步创建的目录
WORKDIR /go/src/go-app

#运行govendor sync命令拉取vendor.json的依赖包
RUN govendor sync

#直接运行main.go，当然你也可以用go build命令编译成二进制文件再运行
CMD go run main.go
