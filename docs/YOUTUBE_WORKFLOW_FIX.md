# YouTube 工作流问题修复指南

## 问题描述

在添加 YouTube 视频分析工作流时，可能会遇到以下两个问题：

1. **迁移文件执行失败**：`migrations/004_youtube_analysis_workflow.sql` 执行失败（可能已执行过）
2. **创建作业失败**：添加 YouTube 视频时报错 `flow has no tasks`

## 问题原因

`flow has no tasks` 错误发生在 `internal/engine/workflow_engine.go:80`，原因是：
- YouTube 工作流 (flows 表) 已创建
- 但流程任务关联 (flow_tasks 表) 没有正确插入
- 导致创建作业时无法找到关联的任务

## 快速修复

### 方法 1: 使用一键修复脚本（推荐）

```bash
make migrate-fix-youtube
```

这个命令会：
1. ✅ 检查 MySQL 连接
2. ✅ 检查 YouTube 工作流状态
3. ✅ 执行迁移文件（如果需要）
4. ✅ 验证修复结果
5. ✅ 显示工作流任务列表

### 方法 2: 分步执行

#### 步骤 1: 设置开发环境（首次使用）

```bash
# 复制环境配置文件
make dev-setup

# 编辑 .env 文件，设置数据库密码
nano .env
```

#### 步骤 2: 执行所有迁移

```bash
make migrate-up
```

#### 步骤 3: 验证 YouTube 工作流

```bash
make migrate-verify
```

#### 步骤 4: 如果验证失败，运行修复脚本

```bash
make migrate-fix-youtube
```

## 手动修复（高级用户）

如果自动修复失败，可以手动执行以下操作：

### 1. 连接到 MySQL

```bash
# 如果使用 Docker
docker exec -it workflow-mysql mysql -uroot -pworkflow123 workflow

# 如果使用本地 MySQL
mysql -uroot -p workflow
```

### 2. 检查问题

```sql
-- 检查流程是否存在
SELECT id, name FROM flows WHERE name = 'YouTube视频分析流程';

-- 检查流程关联的任务数量
SELECT
    f.id, f.name, COUNT(ft.id) as task_count
FROM flows f
LEFT JOIN flow_tasks ft ON f.id = ft.flow_id
WHERE f.name = 'YouTube视频分析流程'
GROUP BY f.id, f.name;

-- 如果 task_count = 0，说明存在问题
```

### 3. 清理并重新创建

```sql
-- 删除现有的 YouTube 工作流（如果存在）
DELETE FROM flows WHERE name = 'YouTube视频分析流程';

-- 删除相关任务（如果需要）
DELETE FROM tasks WHERE name IN (
    '视频信息提取', '字幕下载', '内容分析',
    '情感分析', '数据审核', '生成报告'
);
```

### 4. 重新执行迁移

```bash
make migrate-file FILE=migrations/004_youtube_analysis_workflow.sql
```

## 验证修复

修复完成后，验证 YouTube 工作流：

```bash
make migrate-verify
```

成功的输出应该类似于：

```
✓ YouTube 工作流验证成功
  Flow ID: 3
  Flow Name: YouTube视频分析流程
  Task Count: 6
```

## 测试创建作业

修复完成后，可以通过 API 创建 YouTube 视频分析作业：

### 启动应用

```bash
make run
```

### 创建作业（使用 curl）

```bash
# 获取 YouTube 工作流 ID
FLOW_ID=$(mysql -h localhost -u root -p -D workflow -s -N -e \
  "SELECT id FROM flows WHERE name = 'YouTube视频分析流程' LIMIT 1;")

# 创建作业
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d "{
    \"flow_id\": ${FLOW_ID},
    \"job_name\": \"YouTube视频分析-测试-$(date +%Y%m%d%H%M%S)\",
    \"created_by\": 1
  }"
```

成功响应示例：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "flow_id": 3,
    "job_name": "YouTube视频分析-测试-20250114120000",
    "status": "pending",
    "created_by": 1,
    "created_at": "2025-01-14T12:00:00Z"
  }
}
```

## YouTube 工作流任务说明

修复后的 YouTube 工作流包含以下 6 个任务：

| 序号 | 任务名称 | 类型 | 是否可选 | 说明 |
|------|----------|------|----------|------|
| 1 | 视频信息提取 | automated | 否 | 从 YouTube 提取视频元数据 |
| 2 | 字幕下载 | automated | 是 | 下载视频字幕（可选） |
| 3 | 内容分析 | automated | 否 | 分析视频内容主题和关键词 |
| 4 | 情感分析 | automated | 是 | 分析视频评论情感（可选） |
| 5 | 数据审核 | manual | 否 | 人工审核分析结果 |
| 6 | 生成报告 | automated | 否 | 生成分析报告并保存 |

## 常见问题

### Q1: 执行 `make migrate-fix-youtube` 时提示无法连接 MySQL

**解决方案：**

1. 检查 MySQL 是否运行：
   ```bash
   # Docker 方式
   docker-compose ps
   docker-compose up -d

   # 本地方式
   sudo systemctl status mysql
   sudo systemctl start mysql
   ```

2. 检查 `.env` 文件配置：
   ```bash
   cat .env | grep DB_
   ```

### Q2: 迁移执行后仍然显示 `task_count = 0`

**可能原因：**
- 任务表 (tasks) 中没有对应的任务数据
- 迁移文件执行不完整

**解决方案：**
```sql
-- 检查任务是否存在
SELECT * FROM tasks WHERE name LIKE '%视频%';

-- 如果没有，手动执行迁移文件
source migrations/004_youtube_analysis_workflow.sql;
```

### Q3: 提示 "Duplicate entry" 错误

这是正常的警告，表示数据已经存在。脚本会自动处理这种情况。

## 相关文件

- 迁移文件：`migrations/004_youtube_analysis_workflow.sql`
- 修复脚本：`scripts/fix-youtube-workflow.sh`
- 迁移脚本：`scripts/migrate.sh`
- 工作流引擎：`internal/engine/workflow_engine.go:80`

## 获取帮助

如果问题仍未解决，请：

1. 查看完整的错误日志
2. 检查 `migrations/004_youtube_analysis_workflow.sql` 内容
3. 手动连接数据库检查数据状态
4. 提供详细的错误信息以便诊断
