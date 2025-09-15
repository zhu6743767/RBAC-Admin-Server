#!/bin/bash

# RBAC管理员服务器 - 日志系统部署脚本

set -e

echo "🚀 开始部署RBAC日志系统..."

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查Docker和Docker Compose
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker未安装，请先安装Docker${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose未安装，请先安装Docker Compose${NC}"
    exit 1
fi

# 创建必要的目录
echo -e "${YELLOW}📁 创建目录结构...${NC}"
mkdir -p logs/{error,operation,access}
mkdir -p logstash/{pipeline,config}
mkdir -p scripts

# 启动ELK栈
echo -e "${YELLOW}🐳 启动ELK栈...${NC}"
docker-compose -f docker-compose.logging.yml up -d elasticsearch kibana

# 等待Elasticsearch启动
echo -e "${YELLOW}⏳ 等待Elasticsearch启动...${NC}"
until curl -s http://localhost:9200/_cluster/health | grep -q '"status":"green"\|"status":"yellow"'; do
    sleep 5
done

echo -e "${GREEN}✅ Elasticsearch已启动${NC}"

# 启动Filebeat
echo -e "${YELLOW}📊 启动Filebeat...${NC}"
docker-compose -f docker-compose.logging.yml up -d filebeat

# 验证服务状态
echo -e "${YELLOW}🔍 验证服务状态...${NC}"
docker-compose -f docker-compose.logging.yml ps

# 输出访问信息
echo -e "${GREEN}🎉 日志系统部署完成！${NC}"
echo ""
echo "📊 Kibana访问地址: http://localhost:5601"
echo "📈 Elasticsearch API: http://localhost:9200"
echo "📁 日志目录: ./logs/"
echo ""
echo "🔧 常用命令:"
echo "  查看日志: docker-compose -f docker-compose.logging.yml logs -f"
echo "  停止服务: docker-compose -f docker-compose.logging.yml down"
echo "  重启服务: docker-compose -f docker-compose.logging.yml restart"
echo ""
echo "📖 下一步:"
echo "  1. 启动RBAC应用: go run main.go"
echo "  2. 访问Kibana创建索引模式: rbac-logs-*"
echo "  3. 在Discover中查看日志"