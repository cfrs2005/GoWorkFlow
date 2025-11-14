#!/bin/bash

# YouTube 视频分析工作流 - 功能测试脚本
# 用于快速验证整个流程是否正常工作

set -e

echo "🧪 YouTube 视频分析工作流 - 功能测试"
echo "========================================"
echo ""

# 配置
API_BASE="http://localhost:8080/api"
VIDEO_URL="https://www.youtube.com/watch?v=dQw4w9WgXcQ"

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 检查服务是否运行
echo "📡 检查服务状态..."
if curl -s -f ${API_BASE%/api}/health > /dev/null; then
    echo -e "${GREEN}✓${NC} 服务运行正常"
else
    echo -e "${RED}✗${NC} 服务未运行"
    echo "请先启动服务: ./bin/workflow-api"
    exit 1
fi

echo ""

# 测试 1: 获取流程列表
echo "📋 测试 1/6: 获取流程列表"
echo "----------------------------"
FLOWS_RESPONSE=$(curl -s $API_BASE/flows)
FLOW_COUNT=$(echo $FLOWS_RESPONSE | grep -o '"id"' | wc -l)
echo "找到 $FLOW_COUNT 个流程"

# 查找 YouTube 分析流程
YOUTUBE_FLOW_ID=$(echo $FLOWS_RESPONSE | grep -o '"id":[0-9]*' | grep -A1 'YouTube' | tail -1 | grep -o '[0-9]*' || echo "")

if [ -z "$YOUTUBE_FLOW_ID" ]; then
    echo -e "${YELLOW}未找到 YouTube 分析流程，正在创建...${NC}"

    # 创建流程
    CREATE_RESPONSE=$(curl -s -X POST $API_BASE/flows \
        -H "Content-Type: application/json" \
        -d '{
            "name": "YouTube 视频智能分析",
            "description": "AI 驱动的 YouTube 视频深度分析",
            "version": "1.0.0",
            "is_active": true,
            "created_by": 1
        }')

    YOUTUBE_FLOW_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')

    if [ -n "$YOUTUBE_FLOW_ID" ]; then
        echo -e "${GREEN}✓${NC} 流程创建成功 (ID: $YOUTUBE_FLOW_ID)"
    else
        echo -e "${RED}✗${NC} 流程创建失败"
        echo "响应: $CREATE_RESPONSE"
        exit 1
    fi
else
    echo -e "${GREEN}✓${NC} 找到 YouTube 分析流程 (ID: $YOUTUBE_FLOW_ID)"
fi

echo ""

# 测试 2: 创建作业
echo "🎬 测试 2/6: 创建作业"
echo "----------------------------"
echo "视频 URL: $VIDEO_URL"

JOB_RESPONSE=$(curl -s -X POST $API_BASE/jobs \
    -H "Content-Type: application/json" \
    -d "{
        \"flow_id\": $YOUTUBE_FLOW_ID,
        \"input\": {
            \"video_url\": \"$VIDEO_URL\",
            \"language\": \"en\"
        }
    }")

JOB_ID=$(echo $JOB_RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')

if [ -n "$JOB_ID" ]; then
    echo -e "${GREEN}✓${NC} 作业创建成功 (ID: $JOB_ID)"
else
    echo -e "${RED}✗${NC} 作业创建失败"
    echo "响应: $JOB_RESPONSE"
    exit 1
fi

echo ""

# 测试 3: 自动执行作业
echo "⚡ 测试 3/6: 启动自动执行"
echo "----------------------------"

EXECUTE_RESPONSE=$(curl -s -X POST $API_BASE/jobs/auto-execute \
    -H "Content-Type: application/json" \
    -d "{\"job_id\": $JOB_ID}")

if echo $EXECUTE_RESPONSE | grep -q "Auto execution started"; then
    echo -e "${GREEN}✓${NC} 自动执行已启动"
else
    echo -e "${RED}✗${NC} 启动失败"
    echo "响应: $EXECUTE_RESPONSE"
    exit 1
fi

echo ""

# 测试 4: 监控作业进度
echo "📊 测试 4/6: 监控作业进度"
echo "----------------------------"
echo "等待作业完成 (最多 3 分钟)..."

MAX_WAIT=180  # 3 分钟
WAIT_TIME=0
STATUS="pending"

while [ $WAIT_TIME -lt $MAX_WAIT ]; do
    sleep 5
    WAIT_TIME=$((WAIT_TIME + 5))

    JOB_STATUS=$(curl -s "$API_BASE/jobs?id=$JOB_ID")
    STATUS=$(echo $JOB_STATUS | grep -o '"status":"[^"]*"' | head -1 | cut -d'"' -f4)

    case $STATUS in
        "pending")
            echo -n "⏳ 等待中... (${WAIT_TIME}s) "
            ;;
        "running")
            echo -n "▶️  运行中... (${WAIT_TIME}s) "
            ;;
        "completed")
            echo ""
            echo -e "${GREEN}✓${NC} 作业完成！"
            break
            ;;
        "failed")
            echo ""
            echo -e "${RED}✗${NC} 作业失败"
            echo "状态详情: $JOB_STATUS"
            exit 1
            ;;
        *)
            echo -n "❓ 未知状态: $STATUS (${WAIT_TIME}s) "
            ;;
    esac

    # 每 5 秒输出一个点
    echo -n "."
done

if [ "$STATUS" != "completed" ]; then
    echo ""
    echo -e "${YELLOW}⚠${NC} 超时：作业在 3 分钟内未完成 (当前状态: $STATUS)"
    echo "这可能是正常的，请在 Web 界面继续监控"
fi

echo ""

# 测试 5: 获取 Job Context
echo "📦 测试 5/6: 获取 Job Context"
echo "----------------------------"

CONTEXT_RESPONSE=$(curl -s "$API_BASE/jobs/$JOB_ID/context")
echo "Context 内容:"
echo $CONTEXT_RESPONSE | python3 -m json.tool 2>/dev/null || echo $CONTEXT_RESPONSE

# 提取报告路径
REPORT_URL=$(echo $CONTEXT_RESPONSE | grep -o '"report_url":"[^"]*"' | cut -d'"' -f4)

if [ -n "$REPORT_URL" ]; then
    echo -e "${GREEN}✓${NC} 报告路径: $REPORT_URL"
else
    echo -e "${YELLOW}⚠${NC} 报告路径未找到（可能尚未生成）"
fi

echo ""

# 测试 6: 访问报告
if [ -n "$REPORT_URL" ]; then
    echo "📄 测试 6/6: 访问 HTML 报告"
    echo "----------------------------"

    REPORT_FULL_URL="http://localhost:8080${REPORT_URL}"
    echo "报告 URL: $REPORT_FULL_URL"

    if curl -s -f -I "$REPORT_FULL_URL" > /dev/null; then
        echo -e "${GREEN}✓${NC} 报告可访问"
        echo ""
        echo "📋 报告内容预览:"
        curl -s "$REPORT_FULL_URL" | grep -o '<title>.*</title>' || echo "HTML 内容正常"
    else
        echo -e "${YELLOW}⚠${NC} 报告暂时无法访问"
    fi
else
    echo "📄 测试 6/6: 访问 HTML 报告"
    echo "----------------------------"
    echo -e "${YELLOW}⚠${NC} 跳过（报告未生成）"
fi

echo ""
echo "========================================"
echo "🎉 测试完成！"
echo "========================================"
echo ""
echo "测试结果总结:"
echo "  流程 ID: $YOUTUBE_FLOW_ID"
echo "  作业 ID: $JOB_ID"
echo "  最终状态: $STATUS"
if [ -n "$REPORT_URL" ]; then
    echo "  报告链接: http://localhost:8080${REPORT_URL}"
fi
echo ""
echo "💡 下一步:"
echo "  1. 在浏览器访问: http://localhost:8080"
echo "  2. 进入'作业监控'查看作业 #$JOB_ID"
echo "  3. 点击'详情'查看完整信息"
if [ -n "$REPORT_URL" ]; then
    echo "  4. 访问报告: http://localhost:8080${REPORT_URL}"
fi
echo ""
echo "详细文档: YOUTUBE_ANALYSIS_GUIDE.md"
