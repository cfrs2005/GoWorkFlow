# 🎬 YouTube 视频智能分析工作流

> **真实可用的 AI 驱动视频分析系统** - 自动提取字幕、AI 深度分析、生成精美报告

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=flat&logo=mysql&logoColor=white)](https://www.mysql.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## ⚡ 快速启动（3 种方式任选）

### 🚀 方式一：一键脚本（推荐）

```bash
# Linux / macOS
chmod +x quickstart.sh && ./quickstart.sh

# Windows
quickstart.bat
```

**2-3 分钟**完成所有设置并启动！

### 🐳 方式二：Docker（最简单）

```bash
docker-compose -f docker-quickstart.yml up -d
```

**1 分钟**启动（包含 MySQL + 应用）！

### 💻 方式三：手动启动

```bash
# 1. 初始化数据库
mysql -u root -p workflow < migrations/001_init_schema.sql
mysql -u root -p workflow < migrations/002_sample_data.sql
mysql -u root -p workflow < migrations/003_add_job_context.sql
mysql -u root -p workflow < migrations/004_youtube_analysis_workflow.sql

# 2. 构建并启动
go build -o bin/workflow-api cmd/workflow-api/main.go
./bin/workflow-api
```

📖 **详细说明**: [QUICKSTART_YOUTUBE.md](QUICKSTART_YOUTUBE.md)

---

## 🎯 功能演示

### 1. 访问界面

```
http://localhost:8080
```

![Dashboard](https://via.placeholder.com/800x400?text=Dashboard+Screenshot)

### 2. 创建流程

点击 **🔥 YouTube 视频智能分析** 卡片

### 3. 输入视频

```
https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

### 4. 等待分析

⏱️ **1-3 分钟** 自动完成

### 5. 查看报告

📄 精美的 HTML 分析报告

---

## ✨ 核心特性

### 🎥 YouTube 字幕提取
- ✅ 支持多种字幕获取方式（yt-dlp、youtube-transcript-api）
- ✅ 自动降级到模拟数据（演示模式）
- ✅ 支持多种 URL 格式

### 🤖 AI 深度分析
- 📝 **阅读摘要**：300-500字内容概括
- 🗺️ **思维导图**：结构化内容梳理
- ⭐ **重点分析**：5-8个关键要点
- 💡 **个人认知**：深度思考和启发

### 📊 精美报告
- 🎨 紫色主题设计
- 📱 响应式布局（1920px）
- 🖨️ 打印友好
- 💾 可导出为 PDF

---

## 🛠️ 技术栈

### 后端
- **Go 1.21+** - 高性能 Web 服务
- **MySQL 8.0+** - 关系型数据库
- **Clean Architecture** - 分层架构设计

### AI 集成
- **BigModel GLM-4-Air** - 智谱 AI 大模型
- **自定义提示词** - 优化分析效果

### 前端
- **Alpine.js** - 轻量级响应式框架
- **Tailwind CSS** - 实用优先 CSS
- **Chart.js** - 数据可视化

### 工具链
- **yt-dlp** - YouTube 字幕下载（可选）
- **youtube-transcript-api** - Python 字幕 API（可选）

---

## 📦 项目结构

```
GoWorkFlow/
├── cmd/workflow-api/           # 应用入口
├── internal/
│   ├── executor/               # 任务执行器
│   │   ├── executor.go         # 执行器框架
│   │   ├── youtube_asr_executor.go
│   │   ├── bigmodel_executor.go
│   │   └── html_report_executor.go
│   ├── engine/                 # 工作流引擎
│   ├── service/                # 业务服务层
│   └── handler/                # HTTP 处理器
├── web/                        # Web 界面
│   ├── index.html
│   ├── css/style.css
│   └── js/*.js
├── migrations/                 # 数据库迁移
│   └── 004_youtube_analysis_workflow.sql
├── reports/                    # 生成的报告
├── quickstart.sh               # 快速启动脚本
├── test-youtube-workflow.sh    # 功能测试脚本
└── docker-quickstart.yml       # Docker 配置
```

---

## 🔧 配置说明

### 环境变量（.env）

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=workflow

# 服务配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# BigModel API Key（可选，不设置将使用模拟数据）
BIGMODEL_API_KEY=your_api_key_here
```

### 获取 BigModel API Key

1. 访问：https://open.bigmodel.cn/
2. 注册账号
3. 创建 API Key
4. 每天有免费额度

---

## 🧪 测试

### 自动化测试

```bash
chmod +x test-youtube-workflow.sh
./test-youtube-workflow.sh
```

**测试覆盖**:
- ✅ 服务健康检查
- ✅ 流程创建
- ✅ 作业创建和执行
- ✅ 进度监控
- ✅ 报告生成验证

### 手动测试

```bash
# 1. 创建作业
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "flow_id": 1,
    "input": {
      "video_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
    }
  }'

# 2. 自动执行
curl -X POST http://localhost:8080/api/jobs/auto-execute \
  -H "Content-Type: application/json" \
  -d '{"job_id": 1}'

# 3. 查询状态
curl http://localhost:8080/api/jobs?id=1

# 4. 获取报告
curl http://localhost:8080/api/jobs/1/context
```

---

## 📊 工作流程图

```
┌─────────────────────┐
│  用户输入视频 URL    │
└──────────┬──────────┘
           ↓
┌──────────────────────────────┐
│  Task 1: YouTube ASR 获取     │
│  • 提取视频 ID               │
│  • 获取字幕（自动降级）      │
│  • 保存到 Job Context        │
└──────────┬───────────────────┘
           ↓ (1-30秒)
┌──────────────────────────────┐
│  Task 2: BigModel AI 分析    │
│  • 生成阅读摘要              │
│  • 生成思维导图              │
│  • 生成重点分析              │
│  • 生成个人认知              │
└──────────┬───────────────────┘
           ↓ (30-120秒)
┌──────────────────────────────┐
│  Task 3: HTML 报告生成       │
│  • 应用紫色主题              │
│  • Markdown 转 HTML          │
│  • 保存到 ./reports          │
└──────────┬───────────────────┘
           ↓ (1-5秒)
┌──────────────────────────────┐
│  ✅ 分析完成                 │
│  📄 查看精美报告             │
└──────────────────────────────┘
```

---

## 📖 文档

| 文档 | 说明 |
|------|------|
| [QUICKSTART_YOUTUBE.md](QUICKSTART_YOUTUBE.md) | ⚡ 5分钟快速启动指南 |
| [YOUTUBE_ANALYSIS_GUIDE.md](YOUTUBE_ANALYSIS_GUIDE.md) | 📚 完整使用指南 |
| [WEB_IMPLEMENTATION.md](WEB_IMPLEMENTATION.md) | 🎨 Web 界面说明 |
| [CLAUDE.md](CLAUDE.md) | 🏗️ 项目架构文档 |

---

## 🎯 使用场景

### 学习辅助
- 📚 自动总结教学视频
- 🗺️ 生成学习思维导图
- ⭐ 提取重点知识

### 内容创作
- ✍️ 快速了解视频内容
- 📝 生成文章素材
- 💡 获取创作灵感

### 研究分析
- 🔬 批量分析同主题视频
- 📊 对比不同观点
- 📈 提取关键信息

---

## 🚀 性能指标

| 指标 | 数值 |
|------|------|
| 平均处理时间 | 1-3 分钟 |
| 字幕提取 | < 30秒 |
| AI 分析 | 30-120秒 |
| 报告生成 | < 5秒 |
| 并发支持 | 10+ 作业 |

---

## 🔒 安全性

- ✅ 数据库密码加密存储
- ✅ API Key 环境变量配置
- ✅ 非 root 用户运行（Docker）
- ✅ 输入参数验证
- ✅ SQL 注入防护

---

## 🤝 贡献

欢迎贡献代码和反馈！

### 贡献方式
1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 报告问题
- 🐛 Bug 反馈：[GitHub Issues](https://github.com/cfrs2005/GoWorkFlow/issues)
- 💬 功能建议：[GitHub Discussions](https://github.com/cfrs2005/GoWorkFlow/discussions)

---

## 📝 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 🙏 致谢

### 开源项目
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - YouTube 视频下载工具
- [youtube-transcript-api](https://github.com/jdepoix/youtube-transcript-api) - Python 字幕 API
- [Alpine.js](https://alpinejs.dev/) - 轻量级 JS 框架
- [Tailwind CSS](https://tailwindcss.com/) - CSS 框架
- [Chart.js](https://www.chartjs.org/) - 图表库

### AI 服务
- [智谱 AI](https://open.bigmodel.cn/) - BigModel GLM-4-Air

---

## 📞 联系方式

- **作者**: cfrs2005
- **GitHub**: https://github.com/cfrs2005/GoWorkFlow
- **邮箱**: [待补充]

---

## 🎉 开始使用

```bash
# 克隆项目
git clone https://github.com/cfrs2005/GoWorkFlow.git
cd GoWorkFlow

# 快速启动
chmod +x quickstart.sh
./quickstart.sh

# 访问界面
open http://localhost:8080
```

**祝您使用愉快！** 🚀

---

<div align="center">
  <p>如果这个项目对您有帮助，请给一个 ⭐ Star！</p>
  <p>Made with ❤️ by cfrs2005</p>
</div>
