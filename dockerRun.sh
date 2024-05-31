#!/bin/bash

# 使用指定的 Compose 文件构建 Docker 镜像
echo "开始构建 Docker 镜像..."
if sudo docker-compose -f build.yaml build; then
    echo "构建完成，启动服务..."
    if sudo docker-compose -f build.yaml up -d; then
        echo "服务已成功启动"
    else
        echo "服务启动失败"
        exit 1
    fi
else
    echo "构建失败"
    exit 1
fi
