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

# 抓取队列状态（通过 API 端点检查）
echo ""
echo "📦 抓取队列状态:"
QUEUE_JSON=$(curl -s --connect-timeout 3 http://localhost:8080/api/v1/admin/crawl/queue-status \
    -H "Authorization: Bearer $(curl -s -X POST http://localhost:8080/api/v1/auth/login \
        -H 'Content-Type: application/json' \
        -d '{"email":"founder@media-monitor.com","password":"founder123"}' 2>/dev/null | jq -r '.data.access_token // empty')" 2>/dev/null)

if [ -n "$QUEUE_JSON" ] && echo "$QUEUE_JSON" | jq -e '.data' > /dev/null 2>&1; then
    DETAIL_LAG=$(echo "$QUEUE_JSON" | jq -r '.data.detail_queue.lag // 0')
    INGEST_PENDING=$(echo "$QUEUE_JSON" | jq -r '.data.ingest_queue.pending // 0')
    DISPATCH_LEN=$(echo "$QUEUE_JSON" | jq -r '.data.dispatch_queue.length // 0')
    
    if [ "$DETAIL_LAG" -gt 100 ] 2>/dev/null; then
        print_warning "Detail Queue lag: $DETAIL_LAG (>100)"
    else
        print_success "Detail Queue lag: $DETAIL_LAG"
    fi
    
    if [ "$INGEST_PENDING" -gt 200 ] 2>/dev/null; then
        print_warning "Ingest Queue pending: $INGEST_PENDING (>200)"
    else
        print_success "Ingest Queue pending: $INGEST_PENDING"
    fi
    
    print_success "Dispatch Queue length: $DISPATCH_LEN"
else
    print_warning "无法获取队列状态（API 未运行或认证失败）"
fi

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
