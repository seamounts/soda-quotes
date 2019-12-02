毒鸡汤

该项目完全参照 [毒鸡汤](https://github.com/egotong/nows) 开发。源项目用php开发，本项目使用 go 实现了一遍。

### 部署
本地执行：
```sh
go run server.go
```

容器执行:
1. 构建容器镜像
```sh
docker build -t quotes:latest .
```
2. 运行镜像
```sh
docker run -d -p8081:8081 quotes:latest
```   




