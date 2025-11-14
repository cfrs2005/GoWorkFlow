.PHONY: build run test clean migrate-up migrate-down help

# 变量定义
APP_NAME=workflow-api
BUILD_DIR=bin
MAIN_FILE=cmd/workflow-api/main.go

# 默认目标
help:
	@echo "GoWorkFlow Makefile Commands:"
	@echo "  make build         - 编译应用"
	@echo "  make run           - 运行应用"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理构建文件"
	@echo "  make migrate-up    - 运行所有数据库迁移"
	@echo "  make migrate-file FILE=<file> - 运行指定的迁移文件"
	@echo "  make migrate-verify - 验证 YouTube 工作流配置"
	@echo "  make migrate-fix-youtube - 修复 YouTube 工作流问题"
	@echo "  make migrate-down  - 回滚数据库迁移"
	@echo "  make deps          - 安装依赖"
	@echo "  make dev-setup     - 设置开发环境"
	@echo "  make help          - 显示帮助信息"

# 安装依赖
deps:
	go mod download
	go mod tidy

# 编译应用
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# 运行应用
run:
	@echo "Running $(APP_NAME)..."
	go run $(MAIN_FILE)

# 运行测试
test:
	@echo "Running tests..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "Cleaning build files..."
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# 运行数据库迁移
migrate-up:
	@echo "Running database migrations..."
	@bash scripts/migrate.sh up

# 运行单个迁移文件
migrate-file:
	@if [ -z "$(FILE)" ]; then \
		echo "Error: Please specify FILE=migrations/xxx.sql"; \
		exit 1; \
	fi
	@echo "Running migration: $(FILE)..."
	@bash scripts/migrate.sh file $(FILE)

# 验证 YouTube 工作流
migrate-verify:
	@echo "Verifying YouTube workflow..."
	@bash scripts/migrate.sh verify

# 修复 YouTube 工作流问题
migrate-fix-youtube:
	@echo "Fixing YouTube workflow issues..."
	@bash scripts/fix-youtube-workflow.sh

# 数据库迁移变量（可通过环境变量覆盖）
DB_USER ?= root
DB_PASSWORD ?=

# Docker 相关命令
docker-build:
	docker build -t goworkflow:latest .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

# 开发环境设置
dev-setup: deps
	cp config/.env.example .env
	@echo "Development setup complete. Please update .env file with your configuration."
