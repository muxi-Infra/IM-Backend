# 基础镜像，基于golang的alpine镜像构建--编译阶段
FROM golang:alpine AS builder

ENV TZ=Asia/Shanghai
WORKDIR /home/im_backend
COPY . /home/im_backend
ENV GOPROXY https://goproxy.cn,direct
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app main.go

# 使用alpine这个轻量级镜像为基础镜像--运行阶段
FROM alpine AS runner

# 安装时区数据包
RUN apk add --no-cache tzdata

ENV TZ=Asia/Shanghai
WORKDIR /home/im_backend
COPY --from=builder /home/im_backend/app .
EXPOSE 8181
VOLUME ["/data/config"]

# 设置时区环境变量
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

ENTRYPOINT ["./app","-conf","/data/config/config.yaml"]
