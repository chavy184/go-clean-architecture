# 作用：各服务的 Dockerfile：构建应用本体的容器镜像
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o main ./backend/cmd/server/main.go
CMD ["/app/main"]
