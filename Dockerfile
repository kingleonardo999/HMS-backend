# 构建阶段
FROM golang:1.24 AS builder

# 设置工作目录
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码到工作目录
COPY . .

# 编译应用程序
RUN GIN_MODE=release CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -ldflags "-s -w" -trimpath

# 运行阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制构建结果到目标目录
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行应用程序
CMD ["./main"]