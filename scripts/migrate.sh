#!/bin/bash

set -e  # 遇到错误时退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 加载环境变量
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
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

# 检查 MySQL 是否可用
check_mysql() {
    echo -e "${YELLOW}检查 MySQL 连接...${NC}"
    if ! $MYSQL_CMD -e "SELECT 1;" > /dev/null 2>&1; then
        echo -e "${RED}❌ 无法连接到 MySQL。请检查：${NC}"
        echo "  1. MySQL 是否正在运行"
        echo "  2. 数据库配置是否正确 (.env 文件)"
        echo "  3. 用户名密码是否正确"
        echo ""
        echo "当前配置："
        echo "  DB_HOST: $DB_HOST"
        echo "  DB_PORT: $DB_PORT"
        echo "  DB_USER: $DB_USER"
        echo "  DB_NAME: $DB_NAME"
        exit 1
    fi
    echo -e "${GREEN}✓ MySQL 连接成功${NC}"
}

# 检查数据库是否存在
check_database() {
    if ! $MYSQL_CMD -e "USE $DB_NAME;" > /dev/null 2>&1; then
        echo -e "${YELLOW}数据库 '$DB_NAME' 不存在，正在创建...${NC}"
        $MYSQL_CMD -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
        echo -e "${GREEN}✓ 数据库创建成功${NC}"
    fi
}

# 执行单个迁移文件
run_migration() {
    local file=$1
    local filename=$(basename "$file")

    echo -e "${YELLOW}执行: $filename${NC}"

    # 尝试执行迁移
    if $MYSQL_CMD "$DB_NAME" < "$file" 2>&1 | tee /tmp/migration_output.txt; then
        # 检查输出中是否有错误
        if grep -qi "error\|duplicate" /tmp/migration_output.txt; then
            echo -e "${YELLOW}⚠ $filename 执行时有警告（可能已执行过）${NC}"
            return 1
        else
            echo -e "${GREEN}✓ $filename 执行成功${NC}"
            return 0
        fi
    else
        echo -e "${RED}❌ $filename 执行失败${NC}"
        return 1
    fi
}

# 验证 YouTube 工作流
verify_youtube_workflow() {
    echo -e "${YELLOW}验证 YouTube 工作流...${NC}"

    local result=$($MYSQL_CMD "$DB_NAME" -s -N -e "
        SELECT f.id, f.name, COUNT(ft.id) as task_count
        FROM flows f
        LEFT JOIN flow_tasks ft ON f.id = ft.flow_id
        WHERE f.name = 'YouTube 视频智能分析'
        GROUP BY f.id, f.name;
    " 2>&1)

    if [ -z "$result" ]; then
        echo -e "${RED}❌ YouTube 工作流不存在${NC}"
        return 1
    fi

    local task_count=$(echo "$result" | awk '{print $3}')

    if [ "$task_count" -eq 0 ]; then
        echo -e "${RED}❌ YouTube 工作流存在，但没有关联任务 (task_count = 0)${NC}"
        echo -e "${YELLOW}这是导致 'flow has no tasks' 错误的原因${NC}"
        return 1
    fi

    echo -e "${GREEN}✓ YouTube 工作流验证成功${NC}"
    echo "  Flow ID: $(echo "$result" | awk '{print $1}')"
    echo "  Flow Name: YouTube视频分析流程"
    echo "  Task Count: $task_count"
    return 0
}

# 主函数
main() {
    local command=$1

    case "$command" in
        "up")
            check_mysql
            check_database

            echo -e "${GREEN}开始执行数据库迁移...${NC}"
            echo ""

            # 按顺序执行迁移文件
            local migration_files=(
                "migrations/001_init_schema.sql"
                "migrations/002_sample_data.sql"
                "migrations/004_youtube_analysis_workflow.sql"
            )

            local success_count=0
            local skip_count=0

            for file in "${migration_files[@]}"; do
                if [ -f "$file" ]; then
                    if run_migration "$file"; then
                        ((success_count++))
                    else
                        ((skip_count++))
                    fi
                else
                    echo -e "${YELLOW}⚠ 文件不存在: $file${NC}"
                fi
                echo ""
            done

            echo -e "${GREEN}迁移完成！${NC}"
            echo "  成功: $success_count"
            echo "  跳过: $skip_count"
            echo ""

            # 验证 YouTube 工作流
            verify_youtube_workflow
            ;;

        "file")
            local file=$2
            if [ -z "$file" ]; then
                echo -e "${RED}错误: 请指定迁移文件${NC}"
                echo "用法: $0 file <migration_file>"
                exit 1
            fi

            if [ ! -f "$file" ]; then
                echo -e "${RED}错误: 文件不存在: $file${NC}"
                exit 1
            fi

            check_mysql
            check_database
            run_migration "$file"

            # 如果是 YouTube 工作流迁移，验证结果
            if [[ "$file" == *"youtube"* ]]; then
                echo ""
                verify_youtube_workflow
            fi
            ;;

        "verify")
            check_mysql
            verify_youtube_workflow
            ;;

        *)
            echo "用法: $0 {up|file <file>|verify}"
            echo ""
            echo "命令:"
            echo "  up              - 执行所有迁移文件"
            echo "  file <file>     - 执行指定的迁移文件"
            echo "  verify          - 验证 YouTube 工作流配置"
            exit 1
            ;;
    esac
}

main "$@"
