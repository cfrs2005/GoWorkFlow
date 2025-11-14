-- 示例数据
USE workflow;

-- 插入示例任务
INSERT INTO tasks (name, description, task_type, config, is_active) VALUES
('需求评审', '产品需求评审会议', 'approval', '{"approvers": ["PM", "Tech Lead"], "timeout": 3600}', 1),
('技术设计', '完成技术设计文档', 'manual', '{"template": "tech_design_doc.md", "required_sections": ["架构", "接口", "数据库"]}', 1),
('代码开发', '开发功能代码', 'manual', '{"git_branch_pattern": "feature/*"}', 1),
('代码审查', 'Code Review', 'approval', '{"min_approvers": 2, "required_reviewers": ["Senior Developer"]}', 1),
('单元测试', '运行单元测试', 'automated', '{"command": "go test ./...", "coverage_threshold": 80}', 1),
('集成测试', '运行集成测试', 'automated', '{"command": "make integration-test"}', 1),
('部署到测试环境', '部署到测试环境', 'automated', '{"environment": "staging", "rollback_enabled": true}', 1),
('QA测试', 'QA功能测试', 'manual', '{"test_cases": "required"}', 1),
('发布审批', '生产发布审批', 'approval', '{"approvers": ["Tech Lead", "Product Manager"]}', 1),
('生产部署', '部署到生产环境', 'automated', '{"environment": "production", "backup_required": true}', 1);

-- 插入示例流程
INSERT INTO flows (name, description, version, is_active, created_by) VALUES
('标准功能开发流程', '完整的功能开发到上线流程', '1.0.0', 1, 1),
('热修复流程', '紧急Bug修复流程', '1.0.0', 1, 1);

-- 为"标准功能开发流程"添加任务
INSERT INTO flow_tasks (flow_id, task_id, sequence, is_optional, allow_rollback) VALUES
(1, 1, 1, 0, 1),  -- 需求评审
(1, 2, 2, 0, 1),  -- 技术设计
(1, 3, 3, 0, 1),  -- 代码开发
(1, 4, 4, 0, 1),  -- 代码审查
(1, 5, 5, 0, 1),  -- 单元测试
(1, 6, 6, 1, 1),  -- 集成测试（可选）
(1, 7, 7, 0, 1),  -- 部署到测试环境
(1, 8, 8, 0, 1),  -- QA测试
(1, 9, 9, 0, 0),  -- 发布审批（不允许打回）
(1, 10, 10, 0, 0); -- 生产部署（不允许打回）

-- 为"热修复流程"添加任务
INSERT INTO flow_tasks (flow_id, task_id, sequence, is_optional, allow_rollback) VALUES
(2, 3, 1, 0, 1),  -- 代码开发
(2, 4, 2, 0, 1),  -- 代码审查
(2, 5, 3, 0, 1),  -- 单元测试
(2, 9, 4, 0, 0),  -- 发布审批
(2, 10, 5, 0, 0); -- 生产部署
