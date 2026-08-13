#!/bin/bash

# 服务健康检查脚本
# 用法: ./scripts/services/health-check.sh

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

echo "🔍 Checking service health..."
echo ""

HEALTHY_COUNT=0
TOTAL_COUNT=0

# 检查 API 服务
check_service() {
    local name=$1
    local url=$2
    local expected_code=${3:-200}
    
    TOTAL_COUNT=$((TOTAL_COUNT + 1))
    
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 3 "$url" 2>/dev/null || echo "000")
    
    if [ "$HTTP_CODE" = "$expected_code" ]; then
        print_success "$name: Healthy (HTTP $HTTP_CODE)"
        HEALTHY_COUNT=$((HEALTHY_COUNT + 1))
    elif [ "$HTTP_CODE" = "000" ]; then
        print_warning "$name: Not running"
    else
        print_error "$name: Unhealthy (HTTP $HTTP_CODE)"
    fi
}

# 核心服务
check_service "API" "http://localhost:8080/health"
check_service "API Liveness" "http://localhost:8080/health/live"
check_service "API Readiness" "http://localhost:8080/health/ready"

# 监控服务（可选）
check_service "Prometheus" "http://localhost:9090/-/healthy"
check_service "Grafana" "http://localhost:3000/api/health"

# 数据库和 Redis（通过 API 健康检查间接验证）
echo ""

# 总结
echo "======================================"
if [ $HEALTHY_COUNT -eq $TOTAL_COUNT ]; then
    print_success "All services healthy ($HEALTHY_COUNT/$TOTAL_COUNT)"
    exit 0
elif [ $HEALTHY_COUNT -gt 0 ]; then
    print_warning "Some services unhealthy ($HEALTHY_COUNT/$TOTAL_COUNT healthy)"
    exit 1
else
    print_error "No services running ($HEALTHY_COUNT/$TOTAL_COUNT healthy)"
    exit 2
fi
