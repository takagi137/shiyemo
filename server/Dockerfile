# 构建阶段
FROM golang:1.24-alpine AS builder
WORKDIR /app

# 安装依赖
# GOPROXY 做成可覆盖参数：默认走官方源，国内机器可传 https://goproxy.cn,direct
ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=$GOPROXY

COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建二进制
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/main.go

# 运行阶段
FROM alpine:latest
WORKDIR /app

COPY --from=builder /server ./server
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]
