SUBDIR := config
TEMPLATE_FILES := $(wildcard $(SUBDIR)/*.template)
CONFIG_FILES := $(patsubst %.template,%,$(TEMPLATE_FILES))
# 初始化项目配置信息
.PHONY: init
init:
	@echo "检查配置文件..."
	@for file in $(TEMPLATE_FILES); do \
		new_file=$$(echo $$file | sed 's/\.template$$//'); \
		if [ ! -f $$new_file ]; then \
        			cp $$file $$new_file; \
        			rm -f $$file; \
        			echo "生成文件 $$new_file"; \
        		else \
        			echo "文件 $$new_file 已存在，跳过生成"; \
        		fi \
	done
	@if [ ! -f token.json ]; then \
		touch token.json; \
		echo "生成 token.json 文件"; \
	fi

# 本地运行
.PHONY: local
local:init
	@# 检查是否存在 Python 3 环境
	@if command -v python3 &> /dev/null; then \
		echo "Python 3 已安装，运行 pip install..."; \
		if pip3 install -r "$(PWD)/streamlit_app/requirements.txt" -i https://mirrors.aliyun.com/pypi/simple; then \
			echo "依赖项安装成功"; \
		else \
			echo "依赖项安装失败，请检查错误信息"; \
			exit 1; \
		fi; \
	else \
		echo "Python 3 未安装，请先安装 Python 3"; \
		exit 1; \
	fi

	@# 设置环境变量
	@export CGO_ENABLED=1

	@# 格式化 Go 代码
	@if go fmt; then \
		echo "Go 代码格式化成功"; \
	else \
		echo "Go 代码格式化失败，请检查错误信息"; \
		exit 1; \
	fi
	@# 整理 Go 依赖
	@if go mod tidy; then \
		echo "依赖整理成功"; \
	else \
		echo "依赖整理失败，请检查错误信息"; \
		exit 1; \
	fi
	@# 编译 Go 项目
	@if go build .; then \
		echo "项目构建成功"; \
	else \
		echo "项目构建失败，请检查错误信息"; \
		exit 1; \
	fi
	@# 运行 Go 项目
	@if ./wechat-gptbot; then \
		echo "项目运行成功"; \
	else \
		echo "项目运行失败，请检查错误信息"; \
		exit 1; \
	fi




# 本地docker运行
.PHONY: brun
docker:init
	@docker-compose -f build.yaml up --build -d

# 构建依赖的基础镜像并推送到远程仓库
.PHONY: base
base:
	@# 创建一个名为builder的新构建器实例，并设置为当前使用的构建器
	@docker buildx create --name builder --bootstrap --use

	@docker buildx build --platform linux/amd64,linux/arm64 -f baseimage -t 1003941268/python3.11-alpine:latest --push .
	@docker buildx rm


# 打包镜像并推送仓库
.PHONY: push
push:
	@# 检查并删除已存在的构建器实例
	@echo "Checking if buildx builder instance 'pushBuilder' already exists..."
	@if docker buildx inspect pushBuilder > /dev/null 2>&1; then \
		echo "Removing existing buildx builder instance 'pushBuilder'..."; \
		docker buildx rm pushBuilder; \
	fi

	@# 创建一个名为pushBuilder的新构建器实例，并设置为当前使用的构建器
	@echo "Creating and using new buildx builder instance 'pushBuilder'..."
	docker buildx create --name pushBuilder --bootstrap --use

	@# 使用 trap 确保在任务结束后删除构建器
	@echo "Building and pushing Docker image..."
	@trap 'echo "Cleaning up: removing buildx builder instance 'pushBuilder'"; docker buildx rm pushBuilder' EXIT INT TERM; \
		docker buildx build --platform linux/amd64,linux/arm64 -f dockerfile -t 1003941268/python3.11-alpine:latest --push . || exit 1
