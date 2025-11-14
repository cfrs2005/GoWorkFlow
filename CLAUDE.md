# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

GoWorkFlow 是一个基于 Go 和 MySQL 的工作流管理系统，采用 Clean Architecture 模式设计。系统支持任务编排、流程执行、状态追踪和完整的 RESTful API。

## 核心架构

### 分层架构
- **表现层**: `cmd/workflow-api/main.go` + `internal/handler/` - HTTP 请求处理
- **业务逻辑层**: `internal/service/` + `internal/engine/` - 业务规则和工作流引擎
- **数据访问层**: `internal/repository/` + `internal/models/` - 数据持久化
- **公共层**: `pkg/` - 可复用组件

### 核心数据关系
```
Task (任务定义) → Flow (流程定义) → Job (作业实例)
     ↓                ↓                   ↓
   FlowTasks ← → JobTasks ← → JobTaskLogs
```

### 工作流引擎
位置: `internal/engine/workflow_engine.go`
- 核心方法: `CreateJob()`, `StartTask()`, `CompleteTask()`, `FailTask()`, `SkipTask()`, `RollbackTask()`
- 状态机管理: Job 和 JobTask 状态转换
- 事务管理: 确保数据一致性

## 常用开发命令

### 构建和运行
```bash
# 安装依赖
make deps

# 编译应用
make build                    # 输出到 bin/workflow-api

# 运行应用 (开发模式)
make run

# 直接运行
go run cmd/workflow-api/main.go
```

### �和质量检查
```bash
# 运行测试
make test                    # 或 go test -v ./...

# 代码整理
go mod tidy
```

### 数据库操作
```bash
# 运行数据库迁移
make migrate-up DB_USER=root DB_PASSWORD=your_password

# 回滚数据库
make migrate-down

# 手动执行 SQL
mysql -u root -p < migrations/001_init_schema.sql
```

### Docker 开发
```bash
# 构建镜像
make docker-build

# 启动完整环境 (MySQL + 应用)
docker-compose up -d

# 停止服务
docker-compose down
```

### 开发环境设置
```bash
# 初始化开发环境
make dev-setup               # 复制 .env.example 到 .env 并安装依赖
```

## 关键配置

### 环境变量 (.env)
```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=workflow
DB_CHARSET=utf8mb4
```

### 数据库配置
- 数据库: MySQL 8.0+
- 连接池: 最大100连接，空闲10连接，1小时生命周期
- 字符集: utf8mb4
- 核心表: tasks, flows, flow_tasks, jobs, job_tasks, job_task_logs

## 代码组织原则

### 任务类型
- `manual`: 手动任务 - 需要人工干预
- `automated`: 自动化任务 - 可执行脚本或命令
- `approval`: 审批任务 - 需要多方审批

### 状态管理
Job 状态: `pending` → `running` → `completed`/`failed`
JobTask 状态: `pending` → `running` → `completed`/`failed`/`skipped`/`rolled_back`

### API 设计规范
- RESTful 风格，使用标准 HTTP 方法
- 统一响应格式 (`pkg/response/`)
- 错误处理和日志记录
- 支持 CRUD 操作和业务特定接口

## 重要提醒

1. **包管理**: 使用 Go modules (`go mod`)，不是 PDM
2. **数据库操作**: 关键操作使用事务确保一致性
3. **配置文件**: 开发前需要复制 `config/.env.example` 到 `.env`
4. **Docker 优先**: 推荐使用 `docker-compose up -d` 快速启动开发环境
5. **迁移顺序**: 先执行 `001_init_schema.sql` 再执行 `002_sample_data.sql`

## 扩展指南

### 添加新任务类型
1. 在 `internal/models/task.go` 定义类型常量
2. 在工作流引擎实现执行逻辑 (`internal/engine/`)
3. 更新 API 文档和示例

### 数据库变更
1. 创建新的迁移文件在 `migrations/` 目录
2. 更新 `Makefile` 中的 migrate 目标
3. 对应更新 `internal/models/` 结构体

### API 扩展
- 在 `internal/handler/` 添加新的处理器
- 在 `internal/service/` 实现业务逻辑
- 在 `internal/repository/` 添加数据访问方法