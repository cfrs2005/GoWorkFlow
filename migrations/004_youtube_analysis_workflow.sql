-- 004_youtube_analysis_workflow.sql
-- 创建 YouTube 视频分析工作流模板

-- 1. 创建三个任务定义（防止重复插入）
INSERT INTO tasks (name, description, task_type, config, is_active, created_at, updated_at)
SELECT * FROM (
    SELECT
        'YouTube ASR 获取' as name,
        '从 YouTube 视频获取 ASR（自动语音识别）字幕内容' as description,
        'automated' as task_type,
        '{"executor": "youtube_asr", "language": "en"}' as config,
        1 as is_active,
        NOW() as created_at,
        NOW() as updated_at
    UNION ALL
    SELECT
        'BigModel 内容分析',
        '使用智谱 AI BigModel GLM-4-Air 生成阅读摘要、思维导图、重点分析和个人认知',
        'automated',
        '{"executor": "bigmodel_analysis", "model": "glm-4-air"}',
        1,
        NOW(),
        NOW()
    UNION ALL
    SELECT
        'HTML 报告生成',
        '生成精美的 HTML 分析报告',
        'automated',
        '{"executor": "html_report"}',
        1,
        NOW(),
        NOW()
) AS new_tasks
WHERE NOT EXISTS (
    SELECT 1 FROM tasks WHERE name IN ('YouTube ASR 获取', 'BigModel 内容分析', 'HTML 报告生成')
);

-- 2. 创建流程定义（防止重复插入）
INSERT INTO flows (name, description, version, is_active, created_by, created_at, updated_at)
SELECT
    'YouTube 视频智能分析',
    '自动分析 YouTube 视频内容：获取字幕、AI 分析、生成报告。支持阅读摘要、思维导图、重点分析和个人认知四个维度的深度分析。',
    '1.0.0',
    1,
    1,
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM flows WHERE name = 'YouTube 视频智能分析'
);

-- 3. 获取刚创建的 IDs
SET @flow_id = (SELECT id FROM flows WHERE name = 'YouTube 视频智能分析' ORDER BY id DESC LIMIT 1);
SET @task_asr_id = (SELECT id FROM tasks WHERE name = 'YouTube ASR 获取' ORDER BY id DESC LIMIT 1);
SET @task_analysis_id = (SELECT id FROM tasks WHERE name = 'BigModel 内容分析' ORDER BY id DESC LIMIT 1);
SET @task_report_id = (SELECT id FROM tasks WHERE name = 'HTML 报告生成' ORDER BY id DESC LIMIT 1);

-- 4. 创建流程任务关联（防止重复插入）
-- 注意：数据库使用 is_optional 而不是 is_required，0 表示必需，1 表示可选
INSERT INTO flow_tasks (flow_id, task_id, sequence, is_optional, allow_rollback, created_at, updated_at)
SELECT @flow_id, @task_asr_id, 1, 0, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM flow_tasks WHERE flow_id = @flow_id AND sequence = 1
)
UNION ALL
SELECT @flow_id, @task_analysis_id, 2, 0, 1, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM flow_tasks WHERE flow_id = @flow_id AND sequence = 2
)
UNION ALL
SELECT @flow_id, @task_report_id, 3, 0, 0, NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM flow_tasks WHERE flow_id = @flow_id AND sequence = 3
);

-- 5. 验证插入结果
SELECT
    f.id as flow_id,
    f.name as flow_name,
    COUNT(ft.id) as task_count
FROM flows f
LEFT JOIN flow_tasks ft ON f.id = ft.flow_id
WHERE f.name = 'YouTube 视频智能分析'
GROUP BY f.id, f.name;

-- 显示任务详情
SELECT
    ft.sequence AS '序号',
    t.name AS '任务名称',
    t.task_type AS '类型',
    IF(ft.is_optional = 1, '是', '否') AS '可选',
    IF(ft.allow_rollback = 1, '是', '否') AS '可回滚',
    t.config AS '配置'
FROM flow_tasks ft
JOIN tasks t ON ft.task_id = t.id
JOIN flows f ON ft.flow_id = f.id
WHERE f.name = 'YouTube 视频智能分析'
ORDER BY ft.sequence;
