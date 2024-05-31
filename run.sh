#!/bin/bash

# 获取脚本所在的目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"

# 检查是否存在 Python 3 环境
if command -v python3 &> /dev/null; then
    echo "Python 3 已安装，运行 pip install..."
    # 确保使用 pip3 安装，并使用相对路径找到 requirements.txt
    if pip3 install -r "$SCRIPT_DIR/streamlit_app/requirements.txt" -i https://mirrors.aliyun.com/pypi/simple; then
        echo "依赖项安装成功"
    else
        echo "依赖项安装失败，请检查错误信息"
        exit 1
    fi
else
    echo "Python 3 未安装，请先安装 Python 3"
    exit 1
fi

# 设置环境变量
export CGO_ENABLED=1

# 格式化、整理和编译 Go 代码的函数
build_go_project() {
    if go fmt && go mod tidy; then
        echo "Go 代码格式化和依赖整理成功"
    else
        echo "Go 代码格式化或依赖整理失败，请检查错误信息"
        exit 1
    fi

    if go build .; then
        echo "Go 项目构建成功"
    else
        echo "Go 项目构建失败，请检查错误信息"
        exit 1
    fi

    if ./wechat-gptbot; then
        echo "Go 项目运行成功"
    else
        echo "Go 项目运行失败，请检查错误信息"
        exit 1
    fi
}

# 调用构建 Go 项目的函数
build_go_project
