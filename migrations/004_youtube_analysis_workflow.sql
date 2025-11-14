-- 004_youtube_analysis_workflow.sql
-- 创建 YouTube 视频分析工作流模板

-- 1. 创建三个任务定义
INSERT INTO tasks (name, description, task_type, config, is_active, created_at, updated_at) VALUES
(
    'YouTube ASR 获取',
    '从 YouTube 视频获取 ASR（自动语音识别）字幕内容',
    'automated',
    '{"executor": "youtube_asr", "language": "en"}',
    true,
    NOW(),
    NOW()
),
(
    'BigModel 内容分析',
    '使用智谱 AI BigModel GLM-4-Air 生成阅读摘要、思维导图、重点分析和个人认知',
    'automated',
    '{"executor": "bigmodel_analysis", "model": "glm-4-air"}',
    true,
    NOW(),
    NOW()
),
(
    'HTML 报告生成',
    '生成精美的 HTML 分析报告',
    'automated',
    '{"executor": "html_report"}',
    true,
    NOW(),
    NOW()
);

-- 2. 创建流程定义
INSERT INTO flows (name, description, version, is_active, created_by, created_at, updated_at) VALUES
(
    'YouTube 视频智能分析',
    '自动分析 YouTube 视频内容：获取字幕、AI 分析、生成报告。支持阅读摘要、思维导图、重点分析和个人认知四个维度的深度分析。',
    '1.0.0',
    true,
    1,
    NOW(),
    NOW()
);

-- 3. 获取刚创建的 IDs (使用变量，适用于支持的数据库)
SET @flow_id = LAST_INSERT_ID();
SET @task_asr_id = (SELECT id FROM tasks WHERE name = 'YouTube ASR 获取' ORDER BY id DESC LIMIT 1);
SET @task_analysis_id = (SELECT id FROM tasks WHERE name = 'BigModel 内容分析' ORDER BY id DESC LIMIT 1);
SET @task_report_id = (SELECT id FROM tasks WHERE name = 'HTML 报告生成' ORDER BY id DESC LIMIT 1);

-- 4. 创建流程任务关联
INSERT INTO flow_tasks (flow_id, task_id, sequence, is_required, created_at, updated_at) VALUES
(@flow_id, @task_asr_id, 1, true, NOW(), NOW()),
(@flow_id, @task_analysis_id, 2, true, NOW(), NOW()),
(@flow_id, @task_report_id, 3, true, NOW(), NOW());

-- 5. 插入示例说明数据（可选）
-- 这里可以添加一些示例 Job 用于演示
