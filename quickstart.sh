#!/bin/bash

# YouTube 视频分析工作流 - 快速启动脚本
# 适用于: macOS / Linux

set -e  # 遇到错误立即退出

echo "🚀 GoWorkFlow YouTube 视频分析 - 快速启动"
echo "================================================"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查函数
check_command() {
    if command -v $1 &> /dev/null; then
        echo -e "${GREEN}✓${NC} $1 已安装"
        return 0
    else
        echo -e "${YELLOW}✗${NC} $1 未安装"
        return 1
    fi
}

# 步骤 1: 检查依赖
echo ""
echo "📋 步骤 1/5: 检查系统依赖"
echo "----------------------------"

# 必需依赖
if ! check_command go; then
    echo -e "${RED}错误: Go 未安装${NC}"
    echo "请访问 https://golang.org/dl/ 安装 Go 1.16+"
    exit 1
fi

if ! check_command mysql; then
    echo -e "${RED}错误: MySQL 未安装${NC}"
    echo "请安装 MySQL 8.0+ 或使用 Docker: docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password mysql:8"
    exit 1
fi

# 可选依赖（用于真实 YouTube 字幕）
echo ""
echo "可选依赖（用于真实 YouTube 字幕，非必需）:"
check_command yt-dlp || echo "  → 安装: pip install yt-dlp"
check_command python3 && python3 -c "import youtube_transcript_api" 2>/dev/null && echo -e "${GREEN}✓${NC} youtube-transcript-api 已安装" || echo -e "${YELLOW}✗${NC} youtube-transcript-api 未安装 → pip install youtube-transcript-api"

# 步骤 2: 检查配置
echo ""
echo "⚙️  步骤 2/5: 检查配置文件"
echo "----------------------------"

if [ ! -f ".env" ]; then
    echo "创建 .env 文件..."
    cp config/.env.example .env
    echo -e "${GREEN}✓${NC} .env 文件已创建"
else
    echo -e "${GREEN}✓${NC} .env 文件已存在"
fi

# 读取数据库配置
source .env
DB_PASSWORD=${DB_PASSWORD:-"your_password"}

echo ""
echo "当前配置:"
echo "  数据库地址: ${DB_HOST}:${DB_PORT}"
echo "  数据库名称: ${DB_NAME}"
echo "  数据库用户: ${DB_USER}"

# 步骤 3: 初始化数据库
echo ""
echo "🗄️  步骤 3/5: 初始化数据库"
echo "----------------------------"

read -p "数据库密码 [$DB_PASSWORD]: " INPUT_PASSWORD
DB_PASSWORD=${INPUT_PASSWORD:-$DB_PASSWORD}

echo "正在连接数据库..."

# 创建数据库
mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER} -p${DB_PASSWORD} -e "CREATE DATABASE IF NOT EXISTS ${DB_NAME} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null || {
    echo -e "${RED}错误: 无法连接数据库${NC}"
    echo "请检查:"
    echo "  1. MySQL 是否运行: sudo systemctl status mysql"
    echo "  2. 密码是否正确"
    echo "  3. 用户是否有权限"
    exit 1
}

echo -e "${GREEN}✓${NC} 数据库已创建"

# 运行迁移
echo "运行数据库迁移..."
for migration in migrations/*.sql; do
    echo "  → 执行 $(basename $migration)"
    mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER} -p${DB_PASSWORD} ${DB_NAME} < $migration 2>/dev/null || {
        echo -e "${YELLOW}警告: $migration 执行失败（可能已执行过）${NC}"
    }
done

echo -e "${GREEN}✓${NC} 数据库迁移完成"

# 步骤 4: 构建应用
echo ""
echo "🔨 步骤 4/5: 构建应用"
echo "----------------------------"

if [ ! -f "bin/workflow-api" ]; then
    echo "首次构建应用..."
    go build -o bin/workflow-api cmd/workflow-api/main.go
else
    read -p "应用已存在，是否重新构建? [y/N]: " REBUILD
    if [[ $REBUILD =~ ^[Yy]$ ]]; then
        go build -o bin/workflow-api cmd/workflow-api/main.go
    fi
fi

echo -e "${GREEN}✓${NC} 应用构建完成"

# 步骤 5: 创建必要目录
echo ""
echo "📁 步骤 5/5: 创建必要目录"
echo "----------------------------"

mkdir -p reports
mkdir -p web/css web/js
echo -e "${GREEN}✓${NC} 目录已创建"

# 最终检查
echo ""
echo "================================================"
echo -e "${GREEN}✅ 所有准备工作已完成！${NC}"
echo "================================================"

# 提供启动选项
echo ""
echo "选择启动方式:"
echo "  1) 立即启动（前台运行）"
echo "  2) 后台运行"
echo "  3) 仅显示启动命令"
echo ""

read -p "请选择 [1-3]: " CHOICE

case $CHOICE in
    1)
        echo ""
        echo "🚀 启动应用（按 Ctrl+C 停止）..."
        echo ""
        # 更新 .env 中的密码
        sed -i.bak "s/DB_PASSWORD=.*/DB_PASSWORD=${DB_PASSWORD}/" .env
        ./bin/workflow-api
        ;;
    2)
        echo ""
        echo "🚀 后台启动应用..."
        # 更新 .env 中的密码
        sed -i.bak "s/DB_PASSWORD=.*/DB_PASSWORD=${DB_PASSWORD}/" .env
        nohup ./bin/workflow-api > workflow.log 2>&1 &
        PID=$!
        echo -e "${GREEN}✓${NC} 应用已启动 (PID: $PID)"
        echo ""
        echo "访问地址: http://localhost:8080"
        echo "日志文件: workflow.log"
        echo "停止命令: kill $PID"
        echo ""
        echo "查看日志: tail -f workflow.log"
        ;;
    3)
        echo ""
        echo "手动启动命令:"
        echo ""
        echo "  # 更新数据库密码"
        echo "  sed -i 's/DB_PASSWORD=.*/DB_PASSWORD=${DB_PASSWORD}/' .env"
        echo ""
        echo "  # 启动应用"
        echo "  ./bin/workflow-api"
        echo ""
        echo "  # 或后台运行"
        echo "  nohup ./bin/workflow-api > workflow.log 2>&1 &"
        echo ""
        ;;
    *)
        echo "无效选择"
        exit 1
        ;;
esac

echo ""
echo "📖 快速使用指南:"
echo "----------------------------"
echo "1. 访问 Web 界面: http://localhost:8080"
echo "2. 点击 '流程管理' → 红色的 'YouTube 视频智能分析' 卡片"
echo "3. 输入 YouTube 视频地址，例如: https://www.youtube.com/watch?v=dQw4w9WgXcQ"
echo "4. 等待 1-3 分钟，在 '作业监控' 查看进度"
echo "5. 完成后点击 '详情' 查看报告链接"
echo ""
echo "💡 提示:"
echo "  - 首次运行将使用模拟数据（演示效果）"
echo "  - 如需真实分析，设置环境变量: export BIGMODEL_API_KEY=你的密钥"
echo "  - 详细文档: YOUTUBE_ANALYSIS_GUIDE.md"
echo ""
echo -e "${GREEN}祝您使用愉快！🎉${NC}"
