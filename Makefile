# RBAC Admin Server - 简化版Makefile

# Go参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 二进制文件名
BINARY_NAME=rbac_admin_server_simple
BINARY_UNIX=$(BINARY_NAME)_unix

# 编译参数
BUILD_FLAGS=-ldflags="-s -w"

# 默认目标
.PHONY: all
all: build

# 构建简化版服务器
.PHONY: build-simple
build-simple:
	@echo "🔨 正在构建简化版服务器..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) main_simple.go
	@echo "✅ 构建完成: $(BINARY_NAME)"

# 构建完整版服务器
.PHONY: build-full
build-full:
	@echo "🔨 正在构建完整版服务器..."
	$(GOBUILD) $(BUILD_FLAGS) -o rbac_admin_server main.go
	@echo "✅ 构建完成: rbac_admin_server"

# 构建无CGO版本（用于交叉编译）
.PHONY: build-nocgo
build-nocgo:
	@echo "🔨 正在构建无CGO版本..."
	CGO_ENABLED=0 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME)_nocgo main_simple.go
	@echo "✅ 无CGO版本构建完成: $(BINARY_NAME)_nocgo"

# 运行简化版服务器
.PHONY: run-simple
run-simple: build-simple
	@echo "🚀 启动简化版服务器..."
	./$(BINARY_NAME) -env=settings_simple.yaml

# 运行完整版服务器
.PHONY: run-full
run-full: build-full
	@echo "🚀 启动完整版服务器..."
	./rbac_admin_server -env=dev

# 测试
.PHONY: test
test:
	@echo "🧪 运行测试..."
	$(GOTEST) -v ./...

# 清理
.PHONY: clean
clean:
	@echo "🧹 清理构建文件..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)_nocgo
	rm -f rbac_admin_server
	@echo "✅ 清理完成"

# 下载依赖
.PHONY: deps
deps:
	@echo "📦 下载依赖..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ 依赖下载完成"

# 更新依赖
.PHONY: update-deps
update-deps:
	@echo "🔄 更新依赖..."
	$(GOMOD) download
	$(GOMOD) tidy
	$(GOMOD) verify
	@echo "✅ 依赖更新完成"

# 创建数据目录
.PHONY: init
init:
	@echo "📁 初始化数据目录..."
	mkdir -p data
	mkdir -p logs
	@echo "✅ 数据目录创建完成"

# 帮助
.PHONY: help
help:
	@echo "RBAC Admin Server - Makefile帮助"
	@echo "=================================="
	@echo ""
	@echo "可用目标:"
	@echo "  make build-simple    - 构建简化版服务器"
	@echo "  make build-full      - 构建完整版服务器"
	@echo "  make build-nocgo     - 构建无CGO版本（交叉编译）"
	@echo "  make run-simple      - 运行简化版服务器"
	@echo "  make run-full        - 运行完整版服务器"
	@echo "  make test            - 运行测试"
	@echo "  make clean           - 清理构建文件"
	@echo "  make deps            - 下载依赖"
	@echo "  make update-deps     - 更新依赖"
	@echo "  make init            - 初始化数据目录"
	@echo "  make help            - 显示帮助信息"
	@echo ""
	@echo "使用示例:"
	@echo "  make build-simple && make run-simple"
	@echo "  make build-full && make run-full"
	@echo "  make clean && make build-nocgo"