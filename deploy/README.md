
## K8s 环境搭建

```bash
# 部署mysql
kubectl apply -f deploy/mysql/mysql-configmap.yaml 
kubectl apply -f deploy/mysql/deployment-service.yaml

# 部署redis
kubectl apply -f deploy/redis/redis-config.yaml 
kubectl apply -f deploy/redis/deployment-service.yaml
```

## 本地docker环境

### MySQL

搭建本地 MySQL

```bash
# 启动
docker run -it --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:5.6

# 查看
docker ps

# 进入容器
docker exec -it mysql bash

# 登录mysql
mysql -u root -p

# 授权root用户远程登录
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY '123456';
FLUSH PRIVILEGES;

# 退出容器
exit

# 在宿主机登录
mysql -h127.0.0.1 -uroot -p123456

```

### Redis

搭建本地 Redis

```bash
# 启动
docker run -it --name redis -p 6379:6379 redis:6.0
```

### 启动应用

```bash
# 启动应用
docker run --rm --link=mysql --link=redis  -it -p 8080:8080 user-service:v0.0.21

# 检测服务是否正常
➜  ~ curl localhost:8080/ping

# 输出如下说明正常
{"code":0,"message":"Ok","data":{}}
```