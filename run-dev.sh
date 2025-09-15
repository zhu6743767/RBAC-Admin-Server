#!/bin/bash

# 设置文本颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
echo "╔═══════════════════════════════════════╗"
echo "║        🚀 RBAC Admin Server           ║"
echo "║     开发环境启动脚本 (Linux/Mac)      ║"
echo "╚═══════════════════════════════════════╝"
echo -e "${NC}"
echo -e "${GREEN}🌍 启动开发环境...${NC}"
echo -e "${GREEN}📁 使用配置文件: settings_dev.yaml${NC}"
echo -e "${GREEN}🗄️  数据库: SQLite (无需MySQL)${NC}"
echo -e "${GREEN}📊 日志级别: debug${NC}"
echo

# 检查go是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ 错误: Go未安装或未配置环境变量${NC}"
    echo -e "${YELLOW}📥 请访问 https://golang.org/dl/ 下载安装Go${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go环境检查通过${NC}"

# 检查依赖
echo -e "${BLUE}📦 检查项目依赖...${NC}"
if ! go mod tidy > /dev/null 2>&1; then
    echo -e "${RED}❌ 错误: 依赖下载失败${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 依赖检查完成${NC}"

# 检查配置文件
if [[ ! -f "settings_dev.yaml" ]]; then
    echo -e "${RED}❌ 错误: settings_dev.yaml 配置文件不存在${NC}"
    echo -e "${YELLOW}📝 请确保配置文件在正确位置${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 配置文件检查完成${NC}"
echo
echo -e "${GREEN}🚀 正在启动服务器...${NC}"
echo -e "${GREEN}🌐 访问地址: http://localhost:8080${NC}"
echo -e "${GREEN}📚 API文档: http://localhost:8080/swagger/index.html${NC}"
echo
echo -e "${YELLOW}按 Ctrl+C 停止服务器${NC}"
echo -e "${BLUE}═════════════════════════════════════════${NC}"
echo

go run main.go -env=dev

echo
echo -e "${GREEN}服务器已停止${NC}"