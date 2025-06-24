# aq3cms Makefile

# 变量定义
APP_NAME := aq3cms
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date +%Y-%m-%d\ %H:%M:%S)
GO_VERSION := $(shell go version | awk '{print $$3}')
GIT_COMMIT := $(shell git rev-parse HEAD)

# 构建标志
LDFLAGS := -X 'main.Version=$(VERSION)' \
           -X 'main.BuildTime=$(BUILD_TIME)' \
           -X 'main.GoVersion=$(GO_VERSION)' \
           -X 'main.GitCommit=$(GIT_COMMIT)'

# 默认目标
.PHONY: help
help: ## 显示帮助信息
	@echo "aq3cms 开发和部署工具"
	@echo ""
	@echo "可用命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 开发相关
.PHONY: dev
dev: ## 启动开发环境
	docker-compose -f docker-compose.dev.yml up -d

.PHONY: dev-down
dev-down: ## 停止开发环境
	docker-compose -f docker-compose.dev.yml down

.PHONY: dev-logs
dev-logs: ## 查看开发环境日志
	docker-compose -f docker-compose.dev.yml logs -f aq3cms-dev

.PHONY: dev-rebuild
dev-rebuild: ## 重新构建开发环境
	docker-compose -f docker-compose.dev.yml down
	docker-compose -f docker-compose.dev.yml build --no-cache
	docker-compose -f docker-compose.dev.yml up -d

# 构建相关
.PHONY: build
build: ## 构建应用
	CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) ./cmd/server

.PHONY: build-linux
build-linux: ## 构建Linux版本
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-linux ./cmd/server

.PHONY: build-windows
build-windows: ## 构建Windows版本
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-windows.exe ./cmd/server

.PHONY: build-darwin
build-darwin: ## 构建macOS版本
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-darwin ./cmd/server

.PHONY: build-all
build-all: build-linux build-windows build-darwin ## 构建所有平台版本

# 测试相关
.PHONY: test
test: ## 运行测试
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test-race
test-race: ## 运行竞态检测测试
	go test -race -v ./...

.PHONY: benchmark
benchmark: ## 运行基准测试
	go test -bench=. -benchmem ./...

# 代码质量
.PHONY: lint
lint: ## 运行代码检查
	golangci-lint run

.PHONY: fmt
fmt: ## 格式化代码
	go fmt ./...

.PHONY: vet
vet: ## 运行go vet
	go vet ./...

.PHONY: mod-tidy
mod-tidy: ## 整理依赖
	go mod tidy

.PHONY: mod-download
mod-download: ## 下载依赖
	go mod download

# Docker相关
.PHONY: docker-build
docker-build: ## 构建Docker镜像
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

.PHONY: docker-push
docker-push: ## 推送Docker镜像
	docker push $(APP_NAME):$(VERSION)
	docker push $(APP_NAME):latest

.PHONY: docker-run
docker-run: ## 运行Docker容器
	docker run -d --name $(APP_NAME) -p 8080:8080 $(APP_NAME):latest

# 部署相关
.PHONY: deploy
deploy: ## 部署到生产环境
	docker-compose up -d

.PHONY: deploy-down
deploy-down: ## 停止生产环境
	docker-compose down

.PHONY: deploy-logs
deploy-logs: ## 查看生产环境日志
	docker-compose logs -f aq3cms

.PHONY: deploy-rebuild
deploy-rebuild: ## 重新构建生产环境
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# 数据库相关
.PHONY: db-migrate
db-migrate: ## 运行数据库迁移
	./bin/$(APP_NAME) migrate

.PHONY: db-seed
db-seed: ## 填充测试数据
	./bin/$(APP_NAME) seed

.PHONY: db-backup
db-backup: ## 备份数据库
	docker-compose exec mysql mysqldump -u aq3cms -p aq3cms > backup-$(shell date +%Y%m%d_%H%M%S).sql

.PHONY: db-restore
db-restore: ## 恢复数据库 (需要指定文件: make db-restore FILE=backup.sql)
	docker-compose exec -T mysql mysql -u aq3cms -p aq3cms < $(FILE)

# 清理相关
.PHONY: clean
clean: ## 清理构建文件
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html
	rm -f build-errors.log

.PHONY: clean-docker
clean-docker: ## 清理Docker资源
	docker system prune -f
	docker volume prune -f

.PHONY: clean-all
clean-all: clean clean-docker ## 清理所有文件

# 工具安装
.PHONY: install-tools
install-tools: ## 安装开发工具
	go install github.com/cosmtrek/air@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 健康检查
.PHONY: health
health: ## 检查应用健康状态
	curl -f http://localhost:8080/health || exit 1

.PHONY: ready
ready: ## 检查应用就绪状态
	curl -f http://localhost:8080/ready || exit 1

# 版本信息
.PHONY: version
version: ## 显示版本信息
	@echo "App Name: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
