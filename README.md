# AdPushServer 
#### 使用端口：http://localhost:7999，GoVersion：1.10

基于 Go 语言的重构版本

项目主要功能是基于 websocket 的安卓与 server 中间件。为娃娃机搭载的安卓功能板提供通讯解决方案。



## 启动方式
#### 1.godep包管理工具安装

	go get -u -v github.com/tools/godep


#### 2.本地启动方式
       cd src/
       godep go run main.go  -MySql=true (true:线上，false：线下)

## 构建项目以及选用线上线下数据库
    cd src/
	godep go build main.go
	线上数据库
	main -MySql=true
	线下数据库
	main -MySql=false
	默认使用线下数据库


### 测试功能页面地址
	http://localhost:7999


##  CentOS 服务器Docker方式部署项目
	运行./deploy.sh脚本










