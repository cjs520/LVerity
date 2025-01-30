# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache gcc musl-dev git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -o main ./pkg/main.go

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata ttf-dejavu fontconfig

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml ./
COPY --from=builder /app/scripts ./scripts

# 创建字体目录
RUN mkdir -p /usr/share/fonts/custom

# 复制验证码所需的字体文件
COPY fonts/wqy-microhei.ttc /usr/share/fonts/custom/
RUN fc-cache -f -v

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./main"]
