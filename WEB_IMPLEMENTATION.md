# GoWorkFlow 可视化界面实施说明

## 概述

本次更新为 GoWorkFlow 添加了完整的 Web 可视化管理界面，采用**方案三：混合方案（快速启动）**实施。

## 技术栈

### 后端
- **Go 1.x** - 现有后端框架
- **HTTP FileServer** - 静态文件服务
- **Job Context API** - 新增任务间数据共享功能

### 前端
- **HTML5 + CSS3** - 基础结构和样式
- **Alpine.js 3.x** - 轻量级响应式框架（CDN）
- **Tailwind CSS** - 实用优先的 CSS 框架（CDN）
- **Chart.js 4.x** - 数据可视化图表（CDN）

## 核心功能

### 1. Dashboard（仪表盘）
**路径**: `/` 或 `/#dashboard`

**功能**:
- 统计卡片显示：总流程数、总作业数、运行中作业、今日完成数
- 作业状态分布饼图
- 7日趋势折线图
- 最近作业列表

**技术亮点**:
- 自动每30秒刷新数据
- 紫色渐变主题设计
- Chart.js 实时图表

### 2. Flows Management（流程管理）
**路径**: `/#flows`

**功能**:
- 流程列表展示（表格视图）
- 创建/编辑/删除流程
- 快速模板：
  - **JIRA 数据采集流程**：点击即可创建
  - **RobotSN 数据分析流程**：点击即可创建
- 一键运行流程

**操作**:
```javascript
// 创建流程
POST /api/flows
{
  "name": "流程名称",
  "description": "描述",
  "version": "1.0.0",
  "is_active": true
}

// 运行流程
POST /api/jobs { "flow_id": 1, "input": {} }
POST /api/jobs/start { "job_id": 1 }
```

### 3. Tasks Library（任务库）
**路径**: `/#tasks`

**功能**:
- 任务卡片网格展示
- 任务类型筛选（手动/自动化/审批）
- 创建/编辑/删除任务
- JSON 配置编辑器

**任务类型**:
- `manual`: 手动任务
- `automated`: 自动化任务
- `approval`: 审批任务

### 4. Jobs Monitor（作业监控）
**路径**: `/#jobs`

**功能**:
- 实时作业状态监控
- 自动刷新（可开关，默认5秒）
- 作业详情查看：
  - 输入/输出数据
  - 错误信息
  - 执行时间线
- 作业状态筛选

**状态说明**:
- `pending`: 等待中
- `running`: 运行中
- `completed`: 已完成
- `failed`: 失败

## 新增功能

### Job Context API

为支持任务间数据共享，新增了 Job Context 功能：

**数据库表**:
```sql
CREATE TABLE job_context (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    job_id BIGINT NOT NULL,
    context_key VARCHAR(255) NOT NULL,
    context_value TEXT,
    UNIQUE KEY (job_id, context_key)
);
```

**API 端点**:
```bash
# 获取作业上下文
GET /api/jobs/{id}/context
响应: {"key1": "value1", "key2": "value2"}

# 更新作业上下文
PUT /api/jobs/{id}/context
请求: {"session_id": "abc123", "auth_token": "xyz"}
```

**使用场景**:
- JIRA 工作流：任务A获取的链接传递给任务C
- RobotSN 工作流：任务A的登录凭证传递给任务B

## UI 设计规范

### 配色方案（紫色+白色）
```css
--primary: #7C3AED        /* 主紫色 */
--primary-light: #A78BFA  /* 浅紫色 */
--primary-dark: #5B21B6   /* 深紫色 */
--bg-white: #FFFFFF       /* 白色背景 */
--bg-gray: #F9FAFB        /* 浅灰背景 */
```

### 布局规范
- **最大宽度**: 1920px
- **侧边栏宽度**: 240px（固定）
- **主内容区**: 自适应
- **响应式断点**: 1920 / 1440 / 1024

### 组件样式
- **卡片**: 圆角阴影，悬停效果
- **按钮**: 圆角渐变，过渡动画
- **表格**: 斑马纹，悬停高亮
- **状态徽章**: 圆角彩色标签

## 目录结构

```
GoWorkFlow/
├── web/                      # 前端资源目录
│   ├── index.html           # 主页面
│   ├── css/
│   │   └── style.css        # 自定义样式
│   └── js/
│       ├── api.js           # API 客户端
│       ├── app.js           # 主应用逻辑
│       ├── dashboard.js     # 仪表盘页面
│       ├── flows.js         # 流程管理页面
│       ├── tasks.js         # 任务库页面
│       └── jobs.js          # 作业监控页面
│
├── internal/
│   ├── handler/
│   │   ├── router.go                # 更新：添加静态文件服务
│   │   └── job_context_handler.go  # 新增：Job Context 处理器
│   ├── models/
│   │   └── job_context.go          # 新增：Job Context 模型
│   └── repository/
│       └── job_context_repository.go # 新增：Job Context 仓储
│
└── migrations/
    └── 003_add_job_context.sql      # 新增：数据库迁移
```

## 部署和运行

### 1. 数据库迁移

如果数据库已存在，运行新的迁移：

```bash
# 方式一：使用 MySQL 命令行
mysql -u root -p workflow < migrations/003_add_job_context.sql

# 方式二：使用 Makefile（需要配置）
make migrate-up DB_USER=root DB_PASSWORD=your_password
```

### 2. 配置环境变量

确保 `.env` 文件存在（已自动创建）：

```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=workflow
```

### 3. 启动应用

```bash
# 方式一：直接运行编译好的二进制文件
./bin/workflow-api

# 方式二：使用 go run
go run cmd/workflow-api/main.go

# 方式三：使用 Makefile
make run

# 方式四：Docker（推荐生产环境）
docker-compose up -d
```

### 4. 访问界面

打开浏览器访问:
```
http://localhost:8080
```

## API 路由

### 现有 API
```
POST   /api/tasks           - 创建任务
GET    /api/tasks           - 获取任务列表
GET    /api/tasks?id={id}   - 获取单个任务
PUT    /api/tasks           - 更新任务
DELETE /api/tasks?id={id}   - 删除任务

POST   /api/flows           - 创建流程
GET    /api/flows           - 获取流程列表
GET    /api/flows?id={id}   - 获取单个流程
PUT    /api/flows           - 更新流程
DELETE /api/flows?id={id}   - 删除流程

POST   /api/jobs            - 创建作业
GET    /api/jobs            - 获取作业列表
GET    /api/jobs?id={id}    - 获取单个作业
POST   /api/jobs/start      - 启动作业
GET    /api/jobs/next-task?job_id={id} - 获取下一个任务

POST   /api/tasks/start     - 启动任务
POST   /api/tasks/complete  - 完成任务
POST   /api/tasks/fail      - 标记任务失败
POST   /api/tasks/skip      - 跳过任务
POST   /api/tasks/rollback  - 回滚任务
```

### 新增 API
```
GET    /api/jobs/{id}/context  - 获取作业上下文
PUT    /api/jobs/{id}/context  - 更新作业上下文
```

### 静态文件路由
```
GET    /                    - Web 界面（index.html）
GET    /css/*               - CSS 资源
GET    /js/*                - JavaScript 资源
GET    /health              - 健康检查
```

## 下一步工作

以下功能已规划但未在本次实现：

### 高优先级
1. **任务执行器框架**
   - HTTP 请求执行器（网页抓取）
   - 文件下载执行器（批量下载）
   - 数据解析执行器
   - 远程推送执行器

2. **JIRA 工作流实现**
   - 创建完整的5任务流程
   - 实现动态任务生成（多下载链接处理）
   - 配置任务执行器

3. **RobotSN 工作流实现**
   - 创建完整的3任务流程
   - 实现会话管理（登录凭证共享）
   - 报告生成器

### 中优先级
4. **流程设计器**
   - 拖拽式任务编排
   - 可视化流程图
   - 任务依赖关系配置

5. **高级监控**
   - WebSocket 实时更新
   - Gantt 图时间线
   - 详细日志查看器

### 低优先级
6. **用户管理**
   - 用户认证和授权
   - 角色权限控制
   - 审计日志

7. **通知系统**
   - 邮件通知
   - Webhook 集成
   - 钉钉/企业微信通知

## 技术债务

- [ ] 前端状态管理需要更系统化（考虑使用 Vuex/Pinia）
- [ ] API 响应格式需要统一验证
- [ ] 错误处理需要更细粒度
- [ ] 单元测试覆盖率需要提升
- [ ] 性能优化（大数据量下的分页）

## 已知问题

1. **路由冲突风险**: `/api/jobs/` 路由可能与其他 jobs 路由冲突，建议后续重构为 `/api/job-context/`
2. **前端过滤功能**: 任务和作业的筛选功能目前只是 UI，未实现实际过滤逻辑
3. **WebSocket 未实现**: 作业监控目前使用轮询，建议后续升级为 WebSocket
4. **无用户认证**: 当前无任何认证机制，生产环境需要添加

## 贡献指南

### 添加新页面
1. 在 `web/js/` 创建新的 JS 文件
2. 在 `index.html` 添加页面容器和导航链接
3. 在 `app.js` 的 `loadPageContent()` 添加页面加载逻辑
4. 在 `web/js/` 文件末尾引入新的 JS 文件

### 添加新 API
1. 在 `internal/handler/` 创建 handler
2. 在 `internal/repository/` 创建 repository
3. 在 `internal/models/` 创建 model（如需要）
4. 在 `router.go` 注册路由
5. 在 `web/js/api.js` 添加客户端方法

## 许可证

本项目遵循与 GoWorkFlow 主项目相同的许可证。

---

**实施日期**: 2025-11-14
**实施方案**: 方案三 - 混合方案（快速启动）
**实施人员**: Claude Code
**版本**: v1.0.0
