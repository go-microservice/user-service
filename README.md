# user service

## feature

- 注册
- 登录
- 获取用户信息
- 批量获取用户信息
- 更新用户信息
- 更新用户密码

## 开发流程

> 以下命令执行都是在根目录下执行

### 1. 编写 proto 文件

可以使用 `eagle proto add api/user/v1/user.proto` 的方式新建proto文件

执行后，会在对应的目录下生成一个 proto 文件，里面直接写业务定义就可以了。

### 2. 生成对应的pb文件

有两种方式可以生成:

1. `make grpc`
2. `eagle proto client api/user/v1/user.proto`

都会生成两个文件 `api/user/v1/user.pb.go` 和 `api/user/v1/user_grpc.pb.go`

### 3. 生成对应的server文件

`eagle proto server api/user/v1/user.proto`

执行该命令后，会在 `internal/service` 下 多一个 `user_grpc.go` 文件。

### 4. 编写业务逻辑

在 `internal/service/user_svc.go` 里直接写业务逻辑接口。

### 5. 转换service输出到pb

在 `internal/service/user_grpc.go` 中将 `user_grpc.go` 转为pb输出的方式。

### 6. 将业务注册到gRPC server 中

在 `internal/server/grpc.go` 中新增 `v1.RegisterUserServer(grpcServer, service.NewUserService())`

## 运行

在根目录下执行以下命令

`go run main.go`

grpc即可正常启动

## 调试

安装 grpcurl

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

正常返回

```bash
 ➜  ~ grpcurl -plaintext -d '{"email":"12345678@cc.com","username":"admin8","password":"123456"}' localhost:9090 micro.user.v1.UserService/Register
 {
   "username": "admin8"
 }
```

异常返回

```bash
 ➜  ~ grpcurl -plaintext -d '{"email":"1234567@cc.com","username":"admin7","password":"123456"}' localhost:9090 micro.user.v1.UserService/Register
 ERROR:
   Code: Code(20100)
   Message: The user already exists.
   Details:
   1)	{"@type":"type.googleapis.com/micro.user.v1.RegisterRequest","email":"1234567@cc.com","password":"123456","username":"admin7"}
```