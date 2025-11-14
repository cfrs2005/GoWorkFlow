# YouTube 视频智能分析工作流使用指南

## 🎬 功能概述

这是一个**真实可用**的 AI 驱动的 YouTube 视频分析系统，能够：

1. **自动提取字幕**：从 YouTube 视频获取 ASR（自动语音识别）内容
2. **AI 深度分析**：使用智谱 BigModel GLM-4-Air 生成：
   - 📝 阅读摘要（300-500字）
   - 🗺️ 思维导图（结构化内容）
   - ⭐ 重点分析（5-8个关键要点）
   - 💡 个人认知（深度思考和启发）
3. **生成精美报告**：输出紫色主题的 HTML 分析报告

---

## 🚀 快速开始

### 1. 数据库初始化

```bash
# 确保 MySQL 已启动并创建了 workflow 数据库
mysql -u root -p

# 在 MySQL 中执行
CREATE DATABASE IF NOT EXISTS workflow CHARACTER SET utf8mb4;
USE workflow;

# 运行所有迁移脚本
source migrations/001_init_schema.sql;
source migrations/002_sample_data.sql;
source migrations/003_add_job_context.sql;
source migrations/004_youtube_analysis_workflow.sql;
```

### 2. 配置环境变量

编辑 `.env` 文件（已存在）：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=workflow

# 可选：BigModel API Key（如不设置将使用模拟数据）
BIGMODEL_API_KEY=your_api_key_here
```

**获取 BigModel API Key**：
- 访问：https://open.bigmodel.cn/
- 注册并创建 API Key
- 每天有免费额度可用

### 3. 启动应用

```bash
# 方式一：直接运行编译好的二进制文件
./bin/workflow-api

# 方式二：使用 go run
go run cmd/workflow-api/main.go

# 方式三：使用 Docker
docker-compose up -d
```

### 4. 访问 Web 界面

打开浏览器访问：
```
http://localhost:8080
```

---

## 📖 使用流程

### 方法一：通过 Web 界面（推荐）

1. **创建流程**
   - 访问 `http://localhost:8080/#flows`
   - 点击红色的"🔥 YouTube 视频智能分析"卡片
   - 确认创建流程

2. **运行分析**
   - 输入 YouTube 视频地址，例如：
     ```
     https://www.youtube.com/watch?v=dQw4w9WgXcQ
     https://youtu.be/dQw4w9WgXcQ
     dQw4w9WgXcQ （仅视频 ID）
     ```
   - 点击确定启动分析

3. **查看进度**
   - 自动跳转到"作业监控"页面
   - 实时查看任务执行状态
   - 预计耗时：1-3分钟

4. **查看报告**
   - 作业完成后，点击"查看详情"
   - 在输出数据中找到 `report_url` 字段
   - 访问：`http://localhost:8080/reports/youtube_analysis_XXX.html`

### 方法二：通过 API

#### 步骤 1: 创建流程（如已创建可跳过）

```bash
curl -X POST http://localhost:8080/api/flows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "YouTube 视频智能分析",
    "description": "AI 驱动的 YouTube 视频深度分析",
    "version": "1.0.0",
    "is_active": true,
    "created_by": 1
  }'

# 响应示例
{
  "code": 0,
  "data": {
    "id": 1,
    "name": "YouTube 视频智能分析",
    ...
  }
}
```

#### 步骤 2: 创建作业

```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "flow_id": 1,
    "input": {
      "video_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
      "language": "en"
    }
  }'

# 响应示例
{
  "code": 0,
  "data": {
    "id": 10,
    ...
  }
}
```

#### 步骤 3: 自动执行作业

```bash
curl -X POST http://localhost:8080/api/jobs/auto-execute \
  -H "Content-Type: application/json" \
  -d '{
    "job_id": 10
  }'

# 响应示例
{
  "code": 0,
  "data": {
    "message": "Auto execution started",
    "job_id": 10
  }
}
```

#### 步骤 4: 查询作业状态

```bash
curl http://localhost:8080/api/jobs?id=10

# 响应示例
{
  "code": 0,
  "data": {
    "id": 10,
    "status": "completed",  # pending | running | completed | failed
    ...
  }
}
```

#### 步骤 5: 获取作业上下文（包含报告路径）

```bash
curl http://localhost:8080/api/jobs/10/context

# 响应示例
{
  "code": 0,
  "data": {
    "video_id": "dQw4w9WgXcQ",
    "transcript": "Welcome to this video...",
    "summary": "# 视频摘要\n...",
    "mindmap": "# 思维导图\n...",
    "key_points": "# 重点分析\n...",
    "insights": "# 个人认知\n...",
    "report_path": "./reports/youtube_analysis_XXX.html",
    "report_url": "/reports/youtube_analysis_XXX.html"
  }
}
```

#### 步骤 6: 访问报告

```bash
# 浏览器访问
http://localhost:8080/reports/youtube_analysis_XXX.html

# 或下载报告
curl -o report.html http://localhost:8080/reports/youtube_analysis_XXX.html
```

---

## 🔧 高级配置

### 使用真实的 YouTube 字幕

系统支持三种字幕获取方式（按优先级）：

1. **yt-dlp**（推荐）
   ```bash
   # 安装 yt-dlp
   pip install yt-dlp

   # 或使用包管理器
   brew install yt-dlp  # macOS
   apt install yt-dlp   # Ubuntu/Debian
   ```

2. **youtube-transcript-api**（Python）
   ```bash
   pip install youtube-transcript-api
   ```

3. **模拟数据**（默认）
   - 如果上述工具都未安装，将使用内置的模拟数据
   - 适用于演示和测试

### 使用真实的 BigModel API

```bash
# 设置环境变量
export BIGMODEL_API_KEY="your_actual_api_key"

# 或在 .env 文件中添加
echo "BIGMODEL_API_KEY=your_actual_api_key" >> .env

# 重启应用
./bin/workflow-api
```

### 自定义报告输出目录

默认报告保存在 `./reports` 目录，可在代码中修改：

```go
// cmd/workflow-api/main.go
executor.RegisterExecutor(executor.NewHTMLReportExecutor("/custom/path/to/reports"))
```

---

## 📊 工作流执行流程

```
用户输入 YouTube URL
        ↓
┌───────────────────────────────┐
│  Task 1: YouTube ASR 获取      │
│  ├─ 提取视频 ID                │
│  ├─ 尝试 yt-dlp 获取字幕       │
│  ├─ 备选：youtube-transcript  │
│  └─ 默认：模拟数据            │
│  输出: transcript → Job Context│
└───────────────────────────────┘
        ↓
┌───────────────────────────────┐
│  Task 2: BigModel 内容分析     │
│  输入: transcript (from context)│
│  ├─ 生成阅读摘要              │
│  ├─ 生成思维导图              │
│  ├─ 生成重点分析              │
│  └─ 生成个人认知              │
│  输出: 4个分析结果 → Job Context│
└───────────────────────────────┘
        ↓
┌───────────────────────────────┐
│  Task 3: HTML 报告生成         │
│  输入: 所有分析结果 (from context)│
│  ├─ 应用紫色主题模板          │
│  ├─ Markdown 转 HTML           │
│  └─ 保存到 ./reports 目录     │
│  输出: report_path, report_url │
└───────────────────────────────┘
        ↓
    分析完成
```

---

## 🛠️ 故障排查

### 问题 1：数据库迁移失败

```bash
# 症状：创建流程时提示"流程创建失败"
# 原因：未运行 004_youtube_analysis_workflow.sql

# 解决方案
mysql -u root -p workflow < migrations/004_youtube_analysis_workflow.sql
```

### 问题 2：Task 1 执行失败（字幕获取）

```bash
# 症状：Task 1 状态为"failed"，错误信息"yt-dlp not found"
# 原因：未安装字幕提取工具

# 解决方案1：使用模拟数据（无需安装，自动降级）
# 解决方案2：安装 yt-dlp
pip install yt-dlp

# 解决方案3：安装 youtube-transcript-api
pip install youtube-transcript-api
```

### 问题 3：Task 2 返回模拟数据

```bash
# 症状：分析结果是模拟的通用内容
# 原因：未设置 BIGMODEL_API_KEY

# 解决方案
export BIGMODEL_API_KEY="your_actual_api_key"
# 重启应用
```

### 问题 4：无法访问报告

```bash
# 症状：点击报告链接显示 404
# 原因：reports 目录未创建或权限问题

# 解决方案
mkdir -p ./reports
chmod 755 ./reports
```

### 问题 5：作业卡在 "running" 状态

```bash
# 症状：作业一直显示"运行中"，但没有进展
# 原因：任务执行器抛出异常

# 解决方案：查看日志
# 应用会输出详细的执行日志，查找错误信息

# 手动标记任务失败（如需要）
curl -X POST http://localhost:8080/api/tasks/fail \
  -H "Content-Type: application/json" \
  -d '{
    "job_task_id": <task_id>,
    "error_message": "Manual failure"
  }'
```

---

## 📂 文件结构

```
GoWorkFlow/
├── cmd/workflow-api/
│   └── main.go              # 应用入口，注册执行器
├── internal/
│   ├── executor/            # 任务执行器
│   │   ├── executor.go      # 执行器框架
│   │   ├── youtube_asr_executor.go
│   │   ├── bigmodel_executor.go
│   │   └── html_report_executor.go
│   ├── service/
│   │   └── task_executor_service.go  # 任务执行服务
│   └── handler/
│       ├── executor_handler.go        # 执行API
│       └── router.go                  # 路由配置
├── web/
│   ├── index.html           # Web 界面
│   └── js/
│       └── flows.js         # YouTube 流程创建逻辑
├── migrations/
│   └── 004_youtube_analysis_workflow.sql  # 数据库迁移
└── reports/                 # 生成的报告目录
    └── youtube_analysis_*.html
```

---

## 🎨 报告示例

生成的 HTML 报告包含以下部分：

1. **页眉**
   - 视频 ID 和链接
   - 生成时间
   - 紫色渐变主题

2. **阅读摘要**
   - 视频主要内容概括
   - 核心观点提炼
   - 3-5段结构化内容

3. **思维导图**
   - 分层的内容结构
   - 主题、要点、细节
   - Markdown 列表格式

4. **重点分析**
   - 5-8个关键要点
   - 重要性标注（⭐⭐⭐）
   - 每个要点的详细解释

5. **个人认知**
   - 深度思考和启发
   - 应用场景延伸
   - 批判性分析

---

## 🌟 最佳实践

### 1. 选择合适的视频

- ✅ 推荐：教育、技术、演讲类视频
- ✅ 有完整字幕的视频
- ❌ 避免：音乐MV、无对话视频

### 2. 优化 API 使用

- 使用真实的 BigModel API Key 以获得最佳分析效果
- 定期检查 API 额度
- 对于长视频，字幕可能会被截断（当前限制10000字符）

### 3. 报告管理

- 报告文件名包含视频ID和时间戳
- 定期清理旧报告：`find ./reports -name "*.html" -mtime +30 -delete`
- 可以将报告导出为PDF（使用浏览器打印功能）

### 4. 并发控制

- 当前实现是串行执行任务
- 对于大量分析任务，建议：
  - 使用任务队列（如 Redis Queue）
  - 限制并发数量
  - 添加任务优先级

---

## 🔮 未来计划

- [ ] 支持更多视频平台（Bilibili、抖音、快手）
- [ ] 支持中文字幕优先
- [ ] 添加视频摘要截图
- [ ] 支持批量分析
- [ ] 报告导出为 PDF、Markdown
- [ ] 实时分析进度推送（WebSocket）
- [ ] 分析结果缓存（避免重复分析）

---

## 📞 支持

- GitHub Issues: [https://github.com/cfrs2005/GoWorkFlow/issues](https://github.com/cfrs2005/GoWorkFlow/issues)
- 文档: `WEB_IMPLEMENTATION.md`
- 快速启动: `QUICKSTART_WEB.md`

---

**祝您使用愉快！🎉**
