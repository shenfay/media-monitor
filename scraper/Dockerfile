FROM python:3.12-slim

WORKDIR /app

# 依赖安装（利用 Docker 缓存层）
COPY pyproject.toml .
RUN pip install --no-cache-dir .

# 复制源码
COPY src/ src/

# 运行时配置
ENV MM_REDIS_URL=redis://redis:6379/0

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
    CMD python -m scrapers health || exit 1

# 默认启动常驻消费
CMD ["python", "-m", "scrapers", "serve"]
