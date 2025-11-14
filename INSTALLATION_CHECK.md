# ✅ YouTube 分析工作流 - 安装验证清单

## 📋 核心组件检查

### 1. 应用程序
- [✓] `bin/workflow-api` (13.7 MB) - 已构建
- [✓] `cmd/workflow-api/main.go` - 应用入口
- [✓] `.env` - 环境配置文件

### 2. 执行器模块
- [✓] `internal/executor/executor.go` - 执行器框架
- [✓] `internal/executor/youtube_asr_executor.go` - YouTube 字幕提取
- [✓] `internal/executor/bigmodel_executor.go` - AI 内容分析
- [✓] `internal/executor/html_report_executor.go` - HTML 报告生成

### 3. 数据库迁移
- [✓] `migrations/001_init_schema.sql` - 基础表结构
- [✓] `migrations/002_sample_data.sql` - 示例数据
- [✓] `migrations/003_add_job_context.sql` - 任务上下文表
- [✓] `migrations/004_youtube_analysis_workflow.sql` - YouTube 工作流定义

### 4. 快速启动工具
- [✓] `quickstart.sh` (可执行) - Linux/macOS 快速启动
- [✓] `quickstart.bat` - Windows 快速启动
- [✓] `docker-quickstart.yml` - Docker Compose 配置
- [✓] `Dockerfile.quickstart` - Docker 镜像构建
- [✓] `test-youtube-workflow.sh` (可执行) - 自动化测试

### 5. Web 界面
- [✓] `web/index.html` - 主页面
- [✓] `web/css/style.css` - 紫色主题样式
- [✓] `web/js/*.js` - 前端逻辑
- [✓] `reports/` - 报告输出目录

### 6. 文档
- [✓] `README_YOUTUBE.md` - 完整项目文档
- [✓] `QUICKSTART_YOUTUBE.md` - 5分钟快速启动指南
- [✓] `YOUTUBE_ANALYSIS_GUIDE.md` - 详细使用指南
- [✓] `WEB_IMPLEMENTATION.md` - Web 界面说明
- [✓] `CLAUDE.md` - 项目架构文档

---

## 🚀 三种启动方式

### 方式一：一键脚本（推荐）

```bash
# Linux / macOS
./quickstart.sh

# Windows
quickstart.bat
```

**优势**: 自动检查依赖、初始化数据库、构建应用、启动服务

**时间**: 2-3 分钟

---

### 方式二：Docker（最简单）

```bash
docker-compose -f docker-quickstart.yml up -d
```

**优势**: 无需本地安装 Go 或 MySQL，完全容器化

**时间**: 1 分钟（首次需下载镜像）

---

### 方式三：手动启动

```bash
# 1. 初始化数据库（确保 MySQL 已启动）
mysql -u root -p workflow < migrations/001_init_schema.sql
mysql -u root -p workflow < migrations/002_sample_data.sql
mysql -u root -p workflow < migrations/003_add_job_context.sql
mysql -u root -p workflow < migrations/004_youtube_analysis_workflow.sql

# 2. 构建应用
go build -o bin/workflow-api cmd/workflow-api/main.go

# 3. 启动应用
./bin/workflow-api
```

**时间**: 5-10 分钟（取决于手动操作速度）

---

## 🧪 验证安装

### 1. 启动后访问

```
http://localhost:8080
```

应该看到紫色主题的 Dashboard 页面

### 2. 运行自动化测试

```bash
./test-youtube-workflow.sh
```

应该看到：
- ✓ 服务健康检查通过
- ✓ YouTube 工作流创建成功
- ✓ 作业创建成功
- ✓ 任务自动执行
- ✓ HTML 报告生成

### 3. 手动测试 API

```bash
# 健康检查
curl http://localhost:8080/api/health

# 获取流程列表
curl http://localhost:8080/api/flows

# 创建作业
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "flow_id": 1,
    "input": {
      "video_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
    }
  }'
```

---

## 🔧 配置选项

### 必需配置（.env）

```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password  # ⚠️ 修改为实际密码
DB_NAME=workflow
```

### 可选配置

```bash
# BigModel API Key（不设置将使用模拟数据）
BIGMODEL_API_KEY=your_api_key_here
```

**获取方式**: https://open.bigmodel.cn/

---

## 📊 功能验证检查表

- [ ] 应用成功启动，无错误日志
- [ ] 访问 http://localhost:8080 看到 Dashboard
- [ ] 数据库中存在 YouTube 工作流（flow_id=1）
- [ ] 可以创建新作业
- [ ] 自动执行功能正常
- [ ] Task 1 (YouTube ASR) 完成（有字幕数据或模拟数据）
- [ ] Task 2 (BigModel 分析) 完成（有分析结果或模拟数据）
- [ ] Task 3 (HTML 报告) 完成（reports 目录下有 .html 文件）
- [ ] 可以访问生成的报告

---

## 🐛 常见问题

### 问题 1: 数据库连接失败

```
Error: Failed to connect to database
```

**解决方案**:
- 检查 MySQL 是否启动: `mysql -u root -p`
- 验证 .env 中的数据库密码是否正确
- 确保 workflow 数据库已创建: `CREATE DATABASE workflow;`

### 问题 2: 端口已被占用

```
Error: bind: address already in use
```

**解决方案**:
- 修改 .env 中的 SERVER_PORT
- 或停止占用 8080 端口的进程: `lsof -ti:8080 | xargs kill -9`

### 问题 3: 报告生成失败

```
Error: failed to write report
```

**解决方案**:
- 确保 reports 目录存在: `mkdir -p reports`
- 检查目录权限: `chmod 755 reports`

---

## 📈 性能指标

| 指标 | 预期值 | 说明 |
|------|--------|------|
| 应用启动时间 | < 2 秒 | 数据库连接 + 执行器注册 |
| 视频分析时间 | 1-3 分钟 | 字幕提取 + AI 分析 + 报告生成 |
| API 响应时间 | < 100ms | 健康检查、列表查询 |
| 并发作业数 | 10+ | 取决于服务器资源 |

---

## ✅ 安装完成标志

当你看到以下所有输出时，说明系统已成功安装：

```bash
# 1. 应用启动日志
INFO: Starting workflow API server...
INFO: Database connected successfully
INFO: Registering task executors...
INFO: Registered executors: [youtube-asr bigmodel-glm-4-air html-report-generator]
INFO: Server listening on 0.0.0.0:8080

# 2. Web 界面可访问
http://localhost:8080 → 显示紫色主题 Dashboard

# 3. 测试脚本通过
./test-youtube-workflow.sh
✓ 所有测试通过
✓ 报告已生成: reports/youtube_analysis_XXX.html
```

---

## 🎉 下一步

安装验证通过后：

1. **使用 Web 界面**
   - 访问 http://localhost:8080/#flows
   - 点击 "🔥 YouTube 视频智能分析"
   - 输入视频 URL 开始分析

2. **查看文档**
   - `README_YOUTUBE.md` - 完整功能说明
   - `QUICKSTART_YOUTUBE.md` - 快速上手
   - `YOUTUBE_ANALYSIS_GUIDE.md` - 详细指南

3. **自定义开发**
   - 参考 `CLAUDE.md` 了解架构
   - 在 `internal/executor/` 添加新的执行器
   - 扩展 Web 界面功能

---

**祝您使用愉快！** 🚀

如有问题，请查看：
- GitHub Issues: https://github.com/cfrs2005/GoWorkFlow/issues
- 详细文档: `YOUTUBE_ANALYSIS_GUIDE.md`
