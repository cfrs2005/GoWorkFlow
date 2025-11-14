#!/bin/bash

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}YouTube 工作流问题修复脚本${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 加载环境变量
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo -e "${YELLOW}⚠ 未找到 .env 文件，使用默认配置${NC}"
    echo -e "${YELLOW}如需自定义配置，请运行: make dev-setup${NC}"
    echo ""
fi

# 设置默认值
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-}"
DB_NAME="${DB_NAME:-workflow}"

# MySQL 连接参数
MYSQL_CMD="mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER}"
if [ -n "$DB_PASSWORD" ]; then
    MYSQL_CMD="$MYSQL_CMD -p${DB_PASSWORD}"
fi

echo -e "${YELLOW}步骤 1/4: 检查 MySQL 连接...${NC}"
if ! $MYSQL_CMD -e "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${RED}❌ 无法连接到 MySQL${NC}"
    echo ""
    echo "请检查以下事项："
    echo "  1. MySQL 服务是否正在运行"
    echo "  2. 如果使用 Docker，请运行: docker-compose up -d"
    echo "  3. 检查 .env 配置文件中的数据库连接信息"
    echo ""
    echo "当前配置："
    echo "  DB_HOST=$DB_HOST"
    echo "  DB_PORT=$DB_PORT"
    echo "  DB_USER=$DB_USER"
    echo "  DB_NAME=$DB_NAME"
    exit 1
fi
echo -e "${GREEN}✓ MySQL 连接成功${NC}"
echo ""

echo -e "${YELLOW}步骤 2/4: 检查 YouTube 工作流状态...${NC}"

# 检查流程是否存在
FLOW_EXISTS=$($MYSQL_CMD "$DB_NAME" -s -N -e "
    SELECT COUNT(*) FROM flows WHERE name = 'YouTube视频分析流程';
")

if [ "$FLOW_EXISTS" -eq 0 ]; then
    echo -e "${YELLOW}⚠ YouTube 工作流不存在，开始创建...${NC}"
    FLOW_ID=""
else
    FLOW_ID=$($MYSQL_CMD "$DB_NAME" -s -N -e "
        SELECT id FROM flows WHERE name = 'YouTube视频分析流程' LIMIT 1;
    ")
    echo -e "${GREEN}✓ YouTube 工作流已存在 (ID: $FLOW_ID)${NC}"

    # 检查任务关联
    TASK_COUNT=$($MYSQL_CMD "$DB_NAME" -s -N -e "
        SELECT COUNT(*) FROM flow_tasks WHERE flow_id = $FLOW_ID;
    ")

    echo "  关联任务数: $TASK_COUNT"

    if [ "$TASK_COUNT" -eq 0 ]; then
        echo -e "${RED}  ❌ 问题确认: 工作流没有关联任务 (这会导致 'flow has no tasks' 错误)${NC}"
    else
        echo -e "${GREEN}  ✓ 工作流配置正常${NC}"
        echo ""
        echo -e "${GREEN}YouTube 工作流已经正确配置，无需修复！${NC}"
        exit 0
    fi
fi
echo ""

echo -e "${YELLOW}步骤 3/4: 执行 YouTube 工作流迁移...${NC}"

# 执行迁移文件
if [ -f "migrations/004_youtube_analysis_workflow.sql" ]; then
    if $MYSQL_CMD "$DB_NAME" < migrations/004_youtube_analysis_workflow.sql > /tmp/migration_output.txt 2>&1; then
        echo -e "${GREEN}✓ 迁移执行成功${NC}"
    else
        # 检查是否是因为重复执行导致的错误
        if grep -qi "duplicate" /tmp/migration_output.txt; then
            echo -e "${YELLOW}⚠ 迁移文件已执行过，继续验证...${NC}"
        else
            echo -e "${RED}❌ 迁移执行失败${NC}"
            cat /tmp/migration_output.txt
            exit 1
        fi
    fi
else
    echo -e "${RED}❌ 迁移文件不存在: migrations/004_youtube_analysis_workflow.sql${NC}"
    exit 1
fi
echo ""

echo -e "${YELLOW}步骤 4/4: 验证修复结果...${NC}"

# 重新获取流程ID
FLOW_ID=$($MYSQL_CMD "$DB_NAME" -s -N -e "
    SELECT id FROM flows WHERE name = 'YouTube视频分析流程' LIMIT 1;
")

if [ -z "$FLOW_ID" ]; then
    echo -e "${RED}❌ 验证失败: 工作流未创建${NC}"
    exit 1
fi

# 检查任务关联
TASK_COUNT=$($MYSQL_CMD "$DB_NAME" -s -N -e "
    SELECT COUNT(*) FROM flow_tasks WHERE flow_id = $FLOW_ID;
")

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}验证结果${NC}"
echo -e "${BLUE}========================================${NC}"
echo "Flow ID: $FLOW_ID"
echo "Flow Name: YouTube视频分析流程"
echo "Task Count: $TASK_COUNT"
echo ""

if [ "$TASK_COUNT" -eq 0 ]; then
    echo -e "${RED}❌ 修复失败: 工作流仍然没有关联任务${NC}"
    echo ""
    echo "可能的原因："
    echo "  1. 任务表 (tasks) 中没有对应的任务数据"
    echo "  2. 迁移文件执行不完整"
    echo ""
    echo "建议操作："
    echo "  1. 检查 tasks 表: SELECT * FROM tasks WHERE name LIKE '%视频%';"
    echo "  2. 手动检查 migrations/004_youtube_analysis_workflow.sql"
    echo "  3. 如需重新执行，可以先删除现有数据："
    echo "     DELETE FROM flows WHERE name = 'YouTube视频分析流程';"
    exit 1
else
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}✓✓✓ 修复成功！ ✓✓✓${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "YouTube 工作流现在包含 $TASK_COUNT 个任务："

    # 显示任务列表
    $MYSQL_CMD "$DB_NAME" -e "
        SELECT
            ft.sequence AS '序号',
            t.name AS '任务名称',
            t.task_type AS '类型',
            IF(ft.is_optional = 1, '是', '否') AS '可选',
            IF(ft.allow_rollback = 1, '是', '否') AS '可回滚'
        FROM flow_tasks ft
        JOIN tasks t ON ft.task_id = t.id
        WHERE ft.flow_id = $FLOW_ID
        ORDER BY ft.sequence;
    "

    echo ""
    echo -e "${GREEN}现在可以正常创建 YouTube 视频分析作业了！${NC}"
    echo ""
    echo "下一步："
    echo "  1. 启动应用: make run"
    echo "  2. 创建作业: POST /api/v1/jobs"
    echo "     请求体: {\"flow_id\": $FLOW_ID, \"job_name\": \"YouTube视频分析-测试\", \"created_by\": 1}"
fi
