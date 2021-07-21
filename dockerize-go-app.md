### 把你本地的go web应用制成容器镜像

容器化可以基于docker或者containerd，推荐containerd，containerd还没来得及学，本案例采用docker

#### 1. 创建一个简单的go web app

```
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	fmt.Println("Go Web App Started on Port 3001")
	setupRoutes()
	http.ListenAndServe(":3001", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "My Awesome Go App!!!")
}

func userPage(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	rows, _ := db.Query("SELECT * FROM user")
	for rows.Next() {
		var name string
		var age int
		rows.Scan(&name, &age)
		fmt.Fprintf(w, "name=%s, age=%d\n", name, age)
	}
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "guobin:222222@tcp(mysql:3306)/test?charset=utf8")
	if err != nil {
		panic(err)
	}
	return db
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", userPage)
}

```
先不用去管每行代码在做什么，我们只是要学习如何本地化一个go的web app。

#### 2. 定义容器镜像文件`Dockerfile`

```
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
```
注意：这里没有使用go mod，而以另一个go的包管理器`govendor`为例的

#### 3. 构建容器镜像

我们添加一个叫dockerize.sh的shell脚本来构建容器镜像

dockerize.sh:
```
docker build -t qiuguobin/go-web-app:latest .
```

构建命令格式为`docker build -t docker_user/image_name:tag_name .`
- qiuguobin是我在[dockerhub](https://hub.docker.com/)上注册的用户名
- go-web-app是镜像名
- latest是标签名，通常以版本号作为标签名

项目根目录下执行
```
chmod +x ./dockerize.sh
./dockerize.sh
```

构建成功后我们执行`docker images`将会输出

```
REPOSITORY                       TAG           IMAGE ID       CREATED        SIZE
qiuguobin/go-web-app             latest        ee1c59b90288   4 days ago     891MB
```

我们把这个镜像跑起来看看对不对

```
docker run -it -p 3001:3001 qiuguobin/go-web-app
```

输出
```
Go Web App Started on Port 3001
```

再来访问一下
```
http://localhost:3001
```

输出
```
My Awesome Go App!!!
```

说明我们的自定义镜像构建是ok的。

#### 4. 本地镜像上传dockerhub

我们把本地的镜像上传到dockerhub，将来k8s集群才能够从dockerhub下载下来，所以我们先要上传

```
#这里会问你要用户名和密码，看到Login Succeeded表示成功
docker login

#上传dockerhub
docker push guobinqiu/go-web-app
```
你可以在dockerhub的网站上或者通过桌面端查看到你上传的镜像

补充：
这里的镜像是public的，还需要考虑如果把访问权限设置成[private](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/)后对k8s容器拉取时候的差异。另外，如果项目复杂的话，可以考虑用jenkins来一键构建并上传镜像
****