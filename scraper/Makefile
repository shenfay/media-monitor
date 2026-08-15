.PHONY: install dev test serve run docker

install:        ## 安装生产依赖
	pip install .

dev:            ## 安装开发依赖
	pip install -e ".[dev]"

test:           ## 运行测试
	pytest tests/ -v

serve:          ## 本地常驻消费
	python -m crawl serve

run:            ## 本地直跑
	python -m crawl run --source huanqiu --limit 50

build:          ## 构建包
	python -m build

docker:         ## 构建镜像
	docker build -t media-monitor-scraper .

docker-run:     ## 运行容器
	docker run -d --name scraper \
		-e MM_REDIS_URL=redis://host:6379/0 \
		media-monitor-scraper
