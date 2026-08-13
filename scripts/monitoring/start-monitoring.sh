#!/bin/bash

# Prometheus + Grafana 启动脚本
# 用于快速启动监控栈

set -e

echo "🚀 Starting monitoring stack..."
echo ""

# 检查 Prometheus 是否已安装
if ! command -v prometheus &> /dev/null; then
    echo "❌ Prometheus not found. Please install it:"
    echo "   brew install prometheus"
    exit 1
fi

# 检查 Grafana 是否已安装
if ! command -v grafana &> /dev/null && ! command -v grafana-server &> /dev/null; then
    echo "❌ Grafana not found. Please install it:"
    echo "   brew install grafana"
    exit 1
fi

# 创建日志目录
mkdir -p logs

# 1. 启动 Prometheus
echo "📊 Starting Prometheus..."
PROMETHEUS_CONFIG="prometheus.yml"

if [ ! -f "$PROMETHEUS_CONFIG" ]; then
    echo "❌ Prometheus config not found: $PROMETHEUS_CONFIG"
    echo "   Please create prometheus.yml first"
    exit 1
fi

prometheus --config.file="$PROMETHEUS_CONFIG" > logs/prometheus.log 2>&1 &
PROMETHEUS_PID=$!
echo "✅ Prometheus started (PID: $PROMETHEUS_PID)"

# 等待 Prometheus 启动
sleep 2

# 2. 启动 Grafana
echo "📈 Starting Grafana..."
if command -v brew &> /dev/null; then
    brew services start grafana
    echo "✅ Grafana started (via brew services)"
else
    grafana-server --homepath=/usr/local/grafana > logs/grafana.log 2>&1 &
    GRAFANA_PID=$!
    echo "✅ Grafana started (PID: $GRAFANA_PID)"
fi

# 等待 Grafana 启动
sleep 3

# 3. 输出访问信息
echo ""
echo "======================================"
echo "✅ Monitoring stack started!"
echo "======================================"
echo ""
echo "📊 Prometheus:"
echo "   URL: http://localhost:9090"
echo "   Config: $PROMETHEUS_CONFIG"
echo "   Logs: logs/prometheus.log"
echo ""
echo "📈 Grafana:"
echo "   URL: http://localhost:3000"
echo "   Username: admin"
echo "   Password: admin"
echo "   Dashboard: Import grafana-dashboard.json"
echo ""
echo "🎯 Next Steps:"
echo "   1. Visit http://localhost:9090 to verify Prometheus"
echo "   2. Visit http://localhost:3000 and login to Grafana"
echo "   3. Add Prometheus data source (http://localhost:9090)"
echo "   4. Import dashboard from grafana-dashboard.json"
echo "   5. Start your Go application: cd server && go run cmd/api/main.go"
echo ""
echo "🛑 To stop Prometheus: kill $PROMETHEUS_PID"
echo "   To stop Grafana: brew services stop grafana"
echo ""
echo "Press Ctrl+C to stop all services"

# 等待中断
cleanup() {
    echo ""
    echo "🛑 Stopping services..."
    kill $PROMETHEUS_PID 2>/dev/null || true
    if command -v brew &> /dev/null; then
        brew services stop grafana
    fi
    echo "✅ All services stopped"
    exit 0
}

trap cleanup INT TERM
wait $PROMETHEUS_PID
