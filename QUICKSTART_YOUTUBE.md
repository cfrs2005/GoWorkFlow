# ⚡ YouTube 视频分析 - 5分钟快速启动

## 🚀 三种启动方式（任选其一）

---

## 方式一：一键脚本（推荐 ⭐）

### Linux / macOS

```bash
# 1. 给脚本执行权限
chmod +x quickstart.sh

# 2. 运行脚本
./quickstart.sh

# 脚本会自动：
# ✓ 检查依赖（Go、MySQL）
# ✓ 创建配置文件
# ✓ 初始化数据库
# ✓ 构建应用
# ✓ 启动服务
```

### Windows

```batch
REM 双击运行或在命令行执行
quickstart.bat

REM 脚本会引导您完成所有步骤
```

**预计时间**: 2-3 分钟

---

## 方式二：Docker（最简单 🐳）

```bash
# 1. 一键启动（包含 MySQL + 应用）
docker-compose -f docker-quickstart.yml up -d

# 2. 查看日志
docker-compose -f docker-quickstart.yml logs -f goworkflow

# 3. 访问
# http://localhost:8080

# 停止服务
docker-compose -f docker-quickstart.yml down
```

**优点**:
- ✅ 无需安装 Go、MySQL
- ✅ 环境隔离
- ✅ 一键启动/停止

**预计时间**: 1 分钟（首次构建 5 分钟）

---

## 方式三：手动启动（开发者）

### 1. 前置要求

- Go 1.16+
- MySQL 8.0+
- （可选）Python 3 + yt-dlp

### 2. 数据库初始化

```bash
# 登录 MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE workflow CHARACTER SET utf8mb4;
USE workflow;

# 运行迁移
source migrations/001_init_schema.sql;
source migrations/002_sample_data.sql;
source migrations/003_add_job_context.sql;
source migrations/004_youtube_analysis_workflow.sql;
```

### 3. 配置环境

```bash
# 复制配置文件
cp config/.env.example .env

# 编辑 .env，设置数据库密码
vi .env
```

### 4. 构建并启动

```bash
# 构建
go build -o bin/workflow-api cmd/workflow-api/main.go

# 启动
./bin/workflow-api

# 或后台运行
nohup ./bin/workflow-api > workflow.log 2>&1 &
```

**预计时间**: 5 分钟

---

## 📱 使用流程

### 1. 访问 Web 界面

```
http://localhost:8080
```

### 2. 创建 YouTube 分析流程

1. 点击 **流程管理**
2. 点击红色的 **🔥 YouTube 视频智能分析** 卡片
3. 确认创建

### 3. 运行分析

1. 输入 YouTube 视频地址，例如：
   ```
   https://www.youtube.com/watch?v=dQw4w9WgXcQ
   https://youtu.be/dQw4w9WgXcQ
   dQw4w9WgXcQ
   ```

2. 点击确定

### 4. 查看进度

- 自动跳转到 **作业监控**
- 实时查看任务执行状态
- 预计耗时：**1-3 分钟**

### 5. 查看报告

- 点击 **查看详情**
- 在输出数据中找到 `report_url`
- 访问报告：`http://localhost:8080/reports/youtube_analysis_XXX.html`

---

## 🧪 快速测试

运行自动化测试脚本：

```bash
# 赋予执行权限
chmod +x test-youtube-workflow.sh

# 运行测试
./test-youtube-workflow.sh

# 测试内容：
# ✓ 检查服务状态
# ✓ 创建/获取流程
# ✓ 创建作业
# ✓ 自动执行
# ✓ 监控进度
# ✓ 验证报告
```

---

## 📊 完整工作流程图

```
用户输入 YouTube URL
         ↓
┌─────────────────────┐
│  创建作业 (Job)      │
└─────────────────────┘
         ↓
┌─────────────────────┐
│ 自动执行 (Auto Run) │
└─────────────────────┘
         ↓
┌─────────────────────────────────┐
│ Task 1: YouTube ASR 获取         │
│ ├─ 提取视频 ID                   │
│ ├─ 尝试 yt-dlp 获取字幕          │
│ ├─ 备选: youtube-transcript      │
│ └─ 默认: 模拟数据 (演示)         │
│ 输出: transcript → Job Context   │
└─────────────────────────────────┘
         ↓ (1-30秒)
┌─────────────────────────────────┐
│ Task 2: BigModel 内容分析        │
│ 输入: transcript                 │
│ ├─ 生成阅读摘要                  │
│ ├─ 生成思维导图                  │
│ ├─ 生成重点分析                  │
│ └─ 生成个人认知                  │
│ 输出: 4个分析结果 → Job Context  │
└─────────────────────────────────┘
         ↓ (30-120秒)
┌─────────────────────────────────┐
│ Task 3: HTML 报告生成            │
│ 输入: 所有分析结果               │
│ ├─ 应用紫色主题模板              │
│ ├─ Markdown 转 HTML              │
│ └─ 保存到 ./reports 目录         │
│ 输出: report_path, report_url    │
└─────────────────────────────────┘
         ↓ (1-5秒)
    分析完成 ✅
         ↓
    查看精美报告
```

---

## 🎨 报告示例

生成的 HTML 报告包含：

```
┌──────────────────────────────────┐
│  📹 YouTube 视频分析报告          │
│  ================================  │
│  视频: dQw4w9WgXcQ                │
│  生成时间: 2024-11-14 10:30       │
│  ================================  │
│                                   │
│  📝 阅读摘要                      │
│  ├─ 视频主题概述                  │
│  ├─ 核心观点提炼                  │
│  └─ 结构化段落                    │
│                                   │
│  🗺️ 思维导图                      │
│  ├─ 主题层级                      │
│  ├─ 关键要点                      │
│  └─ 具体细节                      │
│                                   │
│  ⭐ 重点分析                      │
│  ├─ 高优先级要点 (⭐⭐⭐)          │
│  ├─ 中优先级要点 (⭐⭐)            │
│  └─ 补充要点 (⭐)                 │
│                                   │
│  💡 个人认知                      │
│  ├─ 内容价值分析                  │
│  ├─ 应用场景延伸                  │
│  ├─ 批判性思考                    │
│  └─ 个人收获                      │
│                                   │
│  Powered by GoWorkFlow            │
└──────────────────────────────────┘
```

---

## 🔧 高级配置

### 使用真实的 BigModel API

```bash
# 设置环境变量
export BIGMODEL_API_KEY="your_actual_api_key"

# 或在 .env 文件中添加
echo "BIGMODEL_API_KEY=your_actual_api_key" >> .env

# 重启应用
```

**获取 API Key**:
- 访问：https://open.bigmodel.cn/
- 注册并创建 API Key
- 每天有免费额度

### 安装 YouTube 字幕工具（可选）

```bash
# 方式 1: yt-dlp（推荐）
pip install yt-dlp

# 方式 2: youtube-transcript-api
pip install youtube-transcript-api

# 验证安装
yt-dlp --version
python3 -c "import youtube_transcript_api; print('OK')"
```

**效果对比**:
- ❌ 无工具：使用模拟数据（演示效果）
- ✅ 有工具：获取真实字幕（真实分析）

---

## ⚠️ 常见问题

### Q1: 数据库连接失败

```bash
# 检查 MySQL 是否运行
sudo systemctl status mysql  # Linux
brew services list | grep mysql  # macOS

# 检查密码
mysql -u root -p

# 检查端口
netstat -an | grep 3306
```

### Q2: 端口 8080 被占用

```bash
# 修改端口
# 编辑 .env
SERVER_PORT=9090

# 重启应用
```

### Q3: 作业卡在 "running" 状态

```bash
# 查看日志
tail -f workflow.log

# 或重启应用
pkill workflow-api
./bin/workflow-api
```

### Q4: 报告显示 404

```bash
# 检查 reports 目录
ls -la reports/

# 检查权限
chmod 755 reports/

# 查看作业上下文
curl http://localhost:8080/api/jobs/1/context
```

---

## 📚 相关文档

- **详细使用指南**: [YOUTUBE_ANALYSIS_GUIDE.md](YOUTUBE_ANALYSIS_GUIDE.md)
- **Web 界面说明**: [WEB_IMPLEMENTATION.md](WEB_IMPLEMENTATION.md)
- **项目架构**: [CLAUDE.md](CLAUDE.md)

---

## 🎯 下一步

### 体验完整功能

1. ✅ 完成基础使用
2. 🔑 设置 BigModel API Key（获取真实分析）
3. 📦 安装 yt-dlp（获取真实字幕）
4. 🎬 分析多个不同类型的视频
5. 📊 对比分析结果

### 自定义开发

1. 📝 查看源码：`internal/executor/`
2. 🔧 创建自定义执行器
3. 🎨 修改报告模板
4. 🔌 集成其他 AI 服务

### 分享和反馈

- ⭐ Star 项目：https://github.com/cfrs2005/GoWorkFlow
- 🐛 报告问题：https://github.com/cfrs2005/GoWorkFlow/issues
- 💬 讨论交流：https://github.com/cfrs2005/GoWorkFlow/discussions

---

**🎉 祝您使用愉快！**

如有问题，请查看 [YOUTUBE_ANALYSIS_GUIDE.md](YOUTUBE_ANALYSIS_GUIDE.md) 的故障排查部分。
