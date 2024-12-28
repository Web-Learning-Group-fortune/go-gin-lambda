# 使用 Go 官方镜像进行构建
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要依赖（如果需要 C 库，才安装这些）
RUN apk add --no-cache gcc musl-dev

# 复制项目文件到镜像中
COPY . .

# 编译 Go 程序
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .

# 使用 AWS Lambda 官方镜像作为运行时
FROM public.ecr.aws/lambda/go:1

# 将编译好的二进制文件复制到运行时镜像
COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}

# 设置 Lambda 启动命令
CMD ["main"]
