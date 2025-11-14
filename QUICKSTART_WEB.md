# GoWorkFlow Web 界面快速启动指南

## 🚀 5分钟快速启动

### 前置条件

- MySQL 8.0+ 已安装并运行
- Go 1.16+ 已安装
- 数据库 `workflow` 已创建

### 步骤 1: 数据库初始化

```bash
# 创建数据库（如果还没创建）
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS workflow CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行初始化脚本
mysql -u root -p workflow < migrations/001_init_schema.sql
mysql -u root -p workflow < migrations/002_sample_data.sql
mysql -u root -p workflow < migrations/003_add_job_context.sql
```

### 步骤 2: 配置环境变量

编辑 `.env` 文件（已自动创建）:

```bash
# 如果需要修改数据库密码
vi .env

# 确保以下配置正确
DB_PASSWORD=your_actual_password
```

### 步骤 3: 启动应用

```bash
# 方式一：使用已编译的二进制文件
./bin/workflow-api

# 方式二：直接运行
go run cmd/workflow-api/main.go
```

### 步骤 4: 访问 Web 界面

打开浏览器访问:

```
http://localhost:8080
```

您应该能看到紫色主题的 GoWorkFlow 管理界面！

---

## 📖 界面功能概览

### 仪表盘 (Dashboard)
- 显示系统统计数据
- 作业状态图表
- 最近作业列表

### 流程管理 (Flows)
- 创建和管理工作流程
- 快速模板：
  - 点击 "JIRA 数据采集流程" 卡片创建 JIRA 工作流
  - 点击 "RobotSN 数据分析流程" 卡片创建 RobotSN 工作流

### 任务库 (Tasks)
- 创建可复用的任务定义
- 支持三种任务类型：
  - 手动任务
  - 自动化任务
  - 审批任务

### 作业监控 (Jobs)
- 实时监控作业执行
- 查看作业详情和日志
- 自动刷新功能

---

## 🎯 快速测试流程

### 测试 1: 创建简单流程

1. 点击 "任务库" 创建一个测试任务:
   ```
   名称: 测试任务A
   类型: 自动化任务
   配置: {}
   ```

2. 点击 "流程管理" → "创建流程":
   ```
   名称: 测试流程
   描述: 我的第一个工作流
   版本: 1.0.0
   ```

3. 点击流程列表中的 "运行" 按钮

4. 在 "作业监控" 中查看执行状态

### 测试 2: 使用快速模板

1. 在 "流程管理" 页面
2. 点击 "JIRA 数据采集流程" 卡片
3. 确认创建
4. 在流程列表中找到新创建的流程

---

## 🔧 常见问题

### Q: 启动失败，提示数据库连接错误？

**A**: 检查以下几点：
```bash
# 1. MySQL 服务是否运行
systemctl status mysql

# 2. .env 文件中的密码是否正确
cat .env | grep DB_PASSWORD

# 3. 数据库是否存在
mysql -u root -p -e "SHOW DATABASES LIKE 'workflow';"
```

### Q: 访问 localhost:8080 显示 404？

**A**: 确保：
```bash
# 1. 应用正在运行
ps aux | grep workflow-api

# 2. 端口 8080 没有被占用
lsof -i :8080

# 3. web 目录存在且包含文件
ls -la web/
```

### Q: 界面显示但没有数据？

**A**: 这是正常的！首次启动数据库为空，您需要：
1. 创建任务
2. 创建流程
3. 创建并运行作业

或者运行示例数据:
```bash
mysql -u root -p workflow < migrations/002_sample_data.sql
```

### Q: 如何修改端口号？

**A**: 编辑 `.env` 文件：
```bash
SERVER_PORT=9090  # 改为你想要的端口
```

---

## 📝 下一步

### 学习资源
- 查看 `WEB_IMPLEMENTATION.md` 了解完整实现细节
- 查看 `CLAUDE.md` 了解项目架构
- 查看 API 文档了解接口定义

### 开发建议
1. 先熟悉界面基本操作
2. 尝试创建简单的工作流
3. 查看作业执行日志
4. 阅读源码理解实现逻辑

### 生产部署
```bash
# 使用 Docker（推荐）
docker-compose up -d

# 或使用 systemd 服务
sudo systemctl enable workflow-api
sudo systemctl start workflow-api
```

---

## 🆘 需要帮助？

- 查看日志: 应用会输出详细的日志信息
- 检查 GitHub Issues
- 阅读项目文档

---

**祝您使用愉快！** 🎉
