.PHONY: all build-frontend build-backend build clean dev help

# 默认目标：构建前端和后端
all: build

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  make build        - 构建前端和后端（默认）"
	@echo "  make build-frontend - 仅构建前端"
	@echo "  make build-backend  - 仅构建后端"
	@echo "  make clean        - 清理构建产物"
	@echo "  make dev          - 开发模式（分别启动前后端）"
	@echo "  make run          - 运行构建后的后端服务"

# 构建前端
build-frontend:
	@echo "正在构建前端..."
	cd web && npm install && npm run build
	@echo "前端构建完成！"

# 构建后端（依赖前端构建）
build-backend: build-frontend
	@echo "正在构建后端..."
	go build -o bin/server ./cmd/main.go
	@echo "后端构建完成！"

# 完整构建
build: build-frontend
	@echo "正在构建后端..."
	go build -o bin/server ./cmd/main.go
	@echo "=========================================="
	@echo "构建完成！"
	@echo "可执行文件位置：bin/server"
	@echo "启动服务：./bin/server"
	@echo "访问地址：http://localhost:8080"
	@echo "=========================================="

# 清理构建产物
clean:
	@echo "清理构建产物..."
	rm -rf web/dist
	rm -rf bin
	rm -rf web/node_modules
	@echo "清理完成！"

# 开发模式
dev:
	@echo "开发模式："
	@echo "1. 在 web 目录下运行 'npm run dev' 启动前端开发服务器"
	@echo "2. 在根目录运行 'go run cmd/main.go' 启动后端服务"
	@echo ""
	@echo "或者使用以下命令同时启动："
	@echo "  (cd web && npm run dev) & go run cmd/main.go"

# 运行构建后的服务
run:
	@echo "启动服务..."
	./bin/server
