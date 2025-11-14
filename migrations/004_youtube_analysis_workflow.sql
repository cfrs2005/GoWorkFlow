-- YouTube 视频分析工作流
USE workflow;

-- 插入 YouTube 分析相关任务
INSERT INTO tasks (name, description, task_type, config, is_active) VALUES
('视频信息提取', '从 YouTube 提取视频元数据（标题、描述、时长等）', 'automated', '{"api": "youtube_data_api", "fields": ["title", "description", "duration", "view_count"]}', 1),
('字幕下载', '下载视频字幕或生成自动字幕', 'automated', '{"languages": ["zh", "en"], "auto_generated": true}', 1),
('内容分析', '分析视频内容主题和关键词', 'automated', '{"nlp_model": "topic_extraction", "keyword_count": 10}', 1),
('情感分析', '分析视频评论情感倾向', 'automated', '{"sentiment_model": "bert", "sample_size": 100}', 1),
('数据审核', '人工审核分析结果', 'manual', '{"checklist": ["数据完整性", "分析准确性"]}', 1),
('生成报告', '生成分析报告并保存', 'automated', '{"format": "pdf", "include_charts": true}', 1);

-- 插入 YouTube 分析流程
INSERT INTO flows (name, description, version, is_active, created_by) VALUES
('YouTube视频分析流程', 'YouTube视频内容分析和报告生成流程', '1.0.0', 1, 1);

-- 获取刚插入的流程ID和任务ID（假设这是第三个流程，任务ID从11开始）
-- 为"YouTube视频分析流程"添加任务
-- 注意：这里使用的ID需要根据实际数据库中的ID进行调整
INSERT INTO flow_tasks (flow_id, task_id, sequence, is_optional, allow_rollback)
SELECT
    (SELECT id FROM flows WHERE name = 'YouTube视频分析流程' LIMIT 1) as flow_id,
    (SELECT id FROM tasks WHERE name = '视频信息提取' LIMIT 1) as task_id,
    1 as sequence,
    0 as is_optional,
    1 as allow_rollback
WHERE NOT EXISTS (
    SELECT 1 FROM flow_tasks ft
    JOIN flows f ON ft.flow_id = f.id
    WHERE f.name = 'YouTube视频分析流程' AND ft.sequence = 1
);

INSERT INTO flow_tasks (flow_id, task_id, sequence, is_optional, allow_rollback)
SELECT
    (SELECT id FROM flows WHERE name = 'YouTube视频分析流程' LIMIT 1) as flow_id,
    (SELECT id FROM tasks WHERE name = '字幕下载' LIMIT 1) as task_id,
    2 as sequence,
    1 as is_optional,
    1 as allow_rollback
WHERE NOT EXISTS (
    SELECT 1 FROM flow_tasks ft
    JOIN flows f ON ft.flow_id = f.id
    WHERE f.name = 'YouTube视频分析流程' AND ft.sequence = 2
);

INSERT INTO flow_tasks (flow_id, task_id, sequence, is_optional, allow_rollback)
SELECT
    (SELECT id FROM flows WHERE name = 'YouTube视频分析流程' LIMIT 1) as flow_id,
    (SELECT id FROM tasks WHERE name = '内容分析' LIMIT 1) as task_id,
    3 as sequence,
    0 as is_optional,
    1 as allow_rollback
WHERE NOT EXISTS (
    SELECT 1 FROM flow_tasks ft
    JOIN flows f ON ft.flow_id = f.id
    WHERE f.name = 'YouTube视频分析流程' AND ft.sequence = 3
);

INSERT INTO flow_tasks (flow_id, task_id, sequence, is_optional, allow_rollback)
SELECT
    (SELECT id FROM flows WHERE name = 'YouTube视频分析流程' LIMIT 1) as flow_id,
    (SELECT id FROM tasks WHERE name = '情感分析' LIMIT 1) as task_id,
    4 as sequence,
    1 as is_optional,
    1 as allow_rollback
WHERE NOT EXISTS (
    SELECT 1 FROM flow_tasks ft
    JOIN flows f ON ft.flow_id = f.id
    WHERE f.name = 'YouTube视频分析流程' AND ft.sequence = 4
);

INSERT INTO flow_tasks (flow_id, task_id, sequence, is_optional, allow_rollback)
SELECT
    (SELECT id FROM flows WHERE name = 'YouTube视频分析流程' LIMIT 1) as flow_id,
    (SELECT id FROM tasks WHERE name = '数据审核' LIMIT 1) as task_id,
    5 as sequence,
    0 as is_optional,
    1 as allow_rollback
WHERE NOT EXISTS (
    SELECT 1 FROM flow_tasks ft
    JOIN flows f ON ft.flow_id = f.id
    WHERE f.name = 'YouTube视频分析流程' AND ft.sequence = 5
);

INSERT INTO flow_tasks (flow_id, task_id, sequence, is_optional, allow_rollback)
SELECT
    (SELECT id FROM flows WHERE name = 'YouTube视频分析流程' LIMIT 1) as flow_id,
    (SELECT id FROM tasks WHERE name = '生成报告' LIMIT 1) as task_id,
    6 as sequence,
    0 as is_optional,
    0 as allow_rollback
WHERE NOT EXISTS (
    SELECT 1 FROM flow_tasks ft
    JOIN flows f ON ft.flow_id = f.id
    WHERE f.name = 'YouTube视频分析流程' AND ft.sequence = 6
);

-- 验证插入结果
SELECT
    f.id as flow_id,
    f.name as flow_name,
    COUNT(ft.id) as task_count
FROM flows f
LEFT JOIN flow_tasks ft ON f.id = ft.flow_id
WHERE f.name = 'YouTube视频分析流程'
GROUP BY f.id, f.name;
