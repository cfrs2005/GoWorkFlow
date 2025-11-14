# GoWorkFlow

一个基于 Go 和 MySQL 的灵活工作流管理系统，支持任务编排、流程执行、任务跳过和打回等功能。

## 特性

- **任务管理**: 定义可重用的任务模板，支持手动、自动化和审批类型
- **流程编排**: 将多个任务组合成完整的业务流程
- **作业执行**: 从流程创建作业实例，按顺序执行任务
- **灵活控制**: 支持任务跳过、打回、失败处理
- **状态追踪**: 完整的任务执行状态和日志记录
- **RESTful API**: 提供完整的 HTTP API 接口

## 核心概念

### Task（任务）
工作流的基本执行单元，定义可重用的任务模板。支持三种类型：
- `manual`: 手动任务
- `automated`: 自动化任务
- `approval`: 审批任务

### Flow（流程）
多个 Task 的有序组合，定义完整的业务流程。每个 Flow 包含：
- 流程名称、描述和版本
- 按顺序排列的任务列表
- 每个任务的配置（是否可选、是否允许打回等）

### Job（作业）
Flow 的一次执行实例。创建 Job 时会：
- 复制 Flow 的所有任务配置
- 为每个任务创建 JobTask 执行记录
- 跟踪整体执行状态

### JobTask（作业任务）
Job 中每个 Task 的具体执行记录，包含：
- 执行状态（pending/running/completed/failed/skipped/rolled_back）
- 执行人、开始时间、完成时间
- 执行结果和错误信息

## 项目结构

```
GoWorkFlow/
├── cmd/
│   └── workflow-api/      # 应用入口
│       └── main.go
├── internal/              # 私有应用代码
│   ├── config/           # 配置管理
│   ├── engine/           # 工作流引擎
│   ├── handler/          # HTTP 处理器
│   ├── middleware/       # 中间件
│   ├── models/           # 数据模型
│   ├── repository/       # 数据访问层
│   └── service/          # 业务逻辑层
├── pkg/                   # 公共库
│   ├── database/         # 数据库连接
│   ├── logger/           # 日志工具
│   └── response/         # HTTP 响应工具
├── migrations/            # 数据库迁移脚本
├── config/               # 配置文件
├── docs/                 # 文档
├── Makefile              # 构建脚本
├── Dockerfile            # Docker 镜像
└── docker-compose.yml    # Docker 编排
```

## 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Make (可选)

### 安装

1. 克隆仓库
```bash
git clone https://github.com/cfrs2005/GoWorkFlow.git
cd GoWorkFlow
```

2. 安装依赖
```bash
go mod download
```

3. 配置环境变量
```bash
cp config/.env.example .env
# 编辑 .env 文件，设置数据库连接信息
```

4. 初始化数据库
```bash
# 方式1：使用 MySQL 命令
mysql -u root -p < migrations/001_init_schema.sql
mysql -u root -p < migrations/002_sample_data.sql

# 方式2：使用 Makefile
make migrate-up DB_USER=root DB_PASSWORD=your_password
```

5. 运行应用
```bash
# 方式1：直接运行
go run cmd/workflow-api/main.go

# 方式2：使用 Makefile
make run

# 方式3：编译后运行
make build
./bin/workflow-api
```

### 使用 Docker

使用 Docker Compose 快速启动整个系统（包括 MySQL）：

```bash
docker-compose up -d
```

访问服务：
- API: http://localhost:8080
- 健康检查: http://localhost:8080/health

## API 文档

### 任务管理

#### 创建任务
```bash
POST /api/tasks
Content-Type: application/json

{
  "name": "代码审查",
  "description": "进行代码审查",
  "task_type": "approval",
  "config": {
    "min_approvers": 2
  }
}
```

#### 获取任务列表
```bash
GET /api/tasks?limit=20&offset=0
```

#### 获取单个任务
```bash
GET /api/tasks?id=1
```

#### 更新任务
```bash
PUT /api/tasks
Content-Type: application/json

{
  "id": 1,
  "name": "更新后的任务名称",
  "description": "更新后的描述"
}
```

#### 删除任务
```bash
DELETE /api/tasks?id=1
```

### 流程管理

#### 创建流程
```bash
POST /api/flows
Content-Type: application/json

{
  "name": "功能开发流程",
  "description": "标准的功能开发流程",
  "version": "1.0.0",
  "task_ids": [1, 2, 3, 4, 5],
  "created_by": 1
}
```

#### 获取流程详情（包含任务）
```bash
GET /api/flows?id=1
```

#### 获取流程列表
```bash
GET /api/flows?limit=20&offset=0
```

### 作业管理

#### 创建作业
```bash
POST /api/jobs
Content-Type: application/json

{
  "flow_id": 1,
  "job_name": "feature-user-auth",
  "created_by": 1
}
```

#### 启动作业
```bash
POST /api/jobs/start
Content-Type: application/json

{
  "job_id": 1
}
```

#### 获取作业详情（包含所有任务）
```bash
GET /api/jobs?id=1
```

#### 获取下一个待执行任务
```bash
GET /api/jobs/next-task?job_id=1
```

### 任务执行

#### 开始执行任务
```bash
POST /api/tasks/start
Content-Type: application/json

{
  "job_task_id": 1,
  "executor_id": 100
}
```

#### 完成任务
```bash
POST /api/tasks/complete
Content-Type: application/json

{
  "job_task_id": 1,
  "result": {
    "status": "success",
    "output": "任务执行成功"
  }
}
```

#### 任务失败
```bash
POST /api/tasks/fail
Content-Type: application/json

{
  "job_task_id": 1,
  "error_message": "执行失败：权限不足"
}
```

#### 跳过任务
```bash
POST /api/tasks/skip
Content-Type: application/json

{
  "job_task_id": 2,
  "operator_id": 100
}
```

#### 打回任务
```bash
POST /api/tasks/rollback
Content-Type: application/json

{
  "job_task_id": 5,
  "operator_id": 100,
  "target_sequence": 2
}
```

## 使用示例

### 完整的工作流执行流程

1. **创建任务定义**
```bash
# 创建需求评审任务
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "需求评审",
    "description": "产品需求评审",
    "task_type": "approval",
    "config": {"approvers": ["PM", "Tech Lead"]}
  }'

# 创建开发任务
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "代码开发",
    "description": "功能开发",
    "task_type": "manual",
    "config": {}
  }'

# 创建测试任务
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试",
    "description": "功能测试",
    "task_type": "automated",
    "config": {"command": "make test"}
  }'
```

2. **创建流程**
```bash
curl -X POST http://localhost:8080/api/flows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "功能开发流程",
    "description": "从需求到发布的完整流程",
    "version": "1.0.0",
    "task_ids": [1, 2, 3],
    "created_by": 1
  }'
```

3. **创建并启动作业**
```bash
# 创建作业
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "flow_id": 1,
    "job_name": "feature-user-login",
    "created_by": 1
  }'

# 启动作业
curl -X POST http://localhost:8080/api/jobs/start \
  -H "Content-Type: application/json" \
  -d '{"job_id": 1}'
```

4. **执行任务**
```bash
# 开始第一个任务
curl -X POST http://localhost:8080/api/tasks/start \
  -H "Content-Type: application/json" \
  -d '{
    "job_task_id": 1,
    "executor_id": 100
  }'

# 完成第一个任务
curl -X POST http://localhost:8080/api/tasks/complete \
  -H "Content-Type: application/json" \
  -d '{
    "job_task_id": 1,
    "result": {"approved": true}
  }'

# 继续执行后续任务...
```

5. **查看作业状态**
```bash
curl http://localhost:8080/api/jobs?id=1
```

## 数据库设计

详细的数据库设计文档请查看 [docs/database_design.md](docs/database_design.md)

主要表结构：
- `tasks`: 任务定义
- `flows`: 流程定义
- `flow_tasks`: 流程任务关联
- `jobs`: 作业实例
- `job_tasks`: 作业任务执行记录
- `job_task_logs`: 作业任务日志

## 配置说明

通过环境变量配置应用：

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| SERVER_HOST | 服务器监听地址 | 0.0.0.0 |
| SERVER_PORT | 服务器端口 | 8080 |
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户 | root |
| DB_PASSWORD | 数据库密码 | (空) |
| DB_NAME | 数据库名称 | workflow |
| DB_CHARSET | 数据库字符集 | utf8mb4 |

## 开发指南

### 添加新的任务类型

1. 在 `internal/models/task.go` 中定义新的任务类型常量
2. 在工作流引擎中实现对应的执行逻辑
3. 更新 API 文档

### 扩展功能

系统采用分层架构，易于扩展：

- **Repository 层**: 数据访问逻辑
- **Service 层**: 业务逻辑
- **Engine 层**: 工作流引擎核心
- **Handler 层**: HTTP 接口

### 运行测试

```bash
go test ./...
```

## 常见问题

### 如何处理并行任务？

当前版本支持顺序执行，并行任务支持可以通过以下方式实现：
1. 在 `flow_tasks` 表中为并行任务设置相同的 `sequence`
2. 在工作流引擎中增加并行执行逻辑

### 如何实现条件分支？

在 `flow_tasks` 表中使用 `condition_config` 字段配置执行条件：
```json
{
  "condition": "previous_task_result.status == 'approved'",
  "skip_if_false": true
}
```

### 如何集成外部系统？

通过 `automated` 类型的任务，在 `config` 中配置：
- Webhook URL
- 命令行脚本
- API 调用配置

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

## 许可证

MIT License

## 联系方式

- GitHub: [@cfrs2005](https://github.com/cfrs2005)
- 项目主页: https://github.com/cfrs2005/GoWorkFlow
