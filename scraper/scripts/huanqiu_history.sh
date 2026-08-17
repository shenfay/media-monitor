#!/bin/bash
# 环球网历史数据并行回刷脚本
#
# 用法:
#   bash run_history.sh [并发数] [起始日期]
#
# 示例:
#   bash run_history.sh              # 默认 5 并发，从 2026-01-01 开始
#   bash run_history.sh 10           # 10 并发，从 2026-01-01 开始
#   bash run_history.sh 5 2025-06-01 # 5 并发，从 2025-06-01 开始
#
# 特点:
#   - 逐子节点查询：每个子节点独立翻页，绕过 offset=10000 限制
#   - 并行执行：多个源同时抓取，大幅提升效率
#   - 精确过滤：只收集 source_name='环球网' 的文章
#   - 自动去重：跨子节点文章自动去重

set -e

# 配置
export MM_DATABASE_DSN="${MM_DATABASE_DSN:-postgresql://shenfay:postgres@localhost:5432/media_monitor?sslmode=disable}"
MAX_PARALLEL=${1:-5}
SINCE_DATE=${2:-2026-01-01}
DELAY=0.5
LOG_DIR="/Users/shenfay/Projects/media-monitor/scraper/logs"

# 所有环球网源
SOURCES=(
  source_hq_world      # 国际
  source_hq_zy         # 要闻
  source_hq_yrd        # 长三角
  source_hq_house      # 房产
  source_hq_media      # 融媒联播
  source_hq_taiwan     # 台湾
  source_hq_lx         # 教育
  source_hq_auto       # 汽车
  source_hq_quality    # 消费
  source_hq_capital    # 产业
  source_hq_biz        # 商业
  source_hq_ent        # 大文娱
  source_hq_cul        # 文化
  source_hq_city       # 城市
  source_hq_hope       # 公益
  source_hq_tech       # 科技
  source_hq_mil        # 军事
  source_hq_energy     # 能源
  source_hq_health     # 健康
  source_hq_v          # 视频
  source_hq_uav        # 无人机
  source_hq_anquan     # 消防
  source_hq_china      # 国内
  source_hq_xy         # 星青年
  source_hq_finance    # 财经
  source_hq_go         # 文旅
  source_hq_women      # 女性
  source_hq_sports     # 体育
  source_hq_art        # 艺术
  source_hq_fashion    # 时尚
  source_hq_oversea    # 海外看中国
  source_hq_lh         # 领航
  source_hq_qinzi      # 亲子
  source_hq_opinion    # 评论
  source_hq_editorial  # 社评
  source_hq_shanrenping # 单仁平
  source_hq_silkroad   # 丝路
  source_hq_russian    # 俄语
  source_hq_english    # 英语
  source_hq_french     # 法语
  source_hq_spanish    # 西班牙语
  source_hq_arabic     # 阿拉伯语
  source_hq_fjxzc      # 新征程
  source_hq_book       # 听书
)

echo "=========================================="
echo "环球网历史数据并行回刷"
echo "=========================================="
echo "源数量: ${#SOURCES[@]}"
echo "并发数: $MAX_PARALLEL"
echo "起始日期: $SINCE_DATE"
echo "请求间隔: ${DELAY}s"
echo "=========================================="

cd /Users/shenfay/Projects/media-monitor/scraper
mkdir -p "$LOG_DIR"

# 单源抓取函数
run_source() {
  local source=$1
  echo "[$(date +%H:%M:%S)] 开始: $source"
  python -m crawl --source "$source" history --since "$SINCE_DATE" --per-node --delay "$DELAY" >> "$LOG_DIR/${source}_history.log" 2>&1
  if [ $? -eq 0 ]; then
    echo "[$(date +%H:%M:%S)] 完成: $source"
  else
    echo "[$(date +%H:%M:%S)] 失败: $source"
  fi
}

export -f run_source
export LOG_DIR DELAY SINCE_DATE

# 并行执行
printf '%s\n' "${SOURCES[@]}" | xargs -P "$MAX_PARALLEL" -I {} bash -c 'run_source "$@"' _ {}

echo ""
echo "=========================================="
echo "全部完成"
echo "日志目录: $LOG_DIR"
echo "=========================================="
