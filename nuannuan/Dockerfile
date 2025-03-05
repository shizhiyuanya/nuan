# 使用官方的Golang镜像作为基础镜像
FROM golang:1.23.4 AS builder

# 设置Go模块代理
ENV GOPROXY=https://goproxy.cn,direct

#工作目录
WORKDIR /opt
ADD .  /opt

#在Docker工作目录下执行命令
RUN go build -o main ./main.go

#暴露端口
EXPOSE 9090

#执行项目的命令
CMD ["/opt/main"]
