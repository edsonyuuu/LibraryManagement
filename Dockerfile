# 使用 golang 的官方镜像作为基础镜像
FROM golang:1.20 AS builder

# 设置 GOPROXY 环境变量
ENV GOPROXY=https://goproxy.cn,direct

# 设置 GOCACHE 环境变量
ENV GOCACHE=/path/to/cache


# 设置工作目录
WORKDIR /app

# 将项目代码复制到容器中
COPY LibraryManagement/ ./LibraryManagement/
COPY go.mod go.sum ./
RUN go mod download

# 进入 LibraryM 目录
WORKDIR /app/LibraryManagement

# 构建项目
RUN go build -o main .

# 使用轻量级的基础镜像
FROM alpine:latest

# 设置工作目录
#WORKDIR /app

# 复制编译好的二进制文件
COPY --from=builder /app/LibraryManagement/main .

# 指定容器启动时执行的命令
CMD ["./main", "-p", "8083"]
