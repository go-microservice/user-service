# account service

## 开发流程

> 以下命令执行都是在根目录下执行

### 1. 编写 proto 文件

可以使用 `eagle proto add api/user/v1/user.proto` 的方式新建proto文件

执行后，会在根目录下生成一个 proto 文件，里面直接写业务逻辑就可以了。

### 2. 生成对应的pb文件

有两种方式可以生成:

1. `make grpc`
2. `eagle proto client api/user/v1/user.proto`

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
