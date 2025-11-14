-- 创建工作流数据库
CREATE DATABASE IF NOT EXISTS workflow CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE workflow;

-- 1. 任务定义表
CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '任务名称',
    description TEXT COMMENT '任务描述',
    task_type VARCHAR(50) NOT NULL COMMENT '任务类型：manual/automated/approval',
    config JSON COMMENT '任务配置（如执行参数、脚本等）',
    is_active TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name),
    INDEX idx_task_type (task_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务定义表';

-- 2. 流程定义表
CREATE TABLE IF NOT EXISTS flows (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '流程名称',
    description TEXT COMMENT '流程描述',
    version VARCHAR(20) DEFAULT '1.0.0' COMMENT '流程版本',
    is_active TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    created_by BIGINT COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name),
    INDEX idx_version (version)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='流程定义表';

-- 3. 流程任务关联表
CREATE TABLE IF NOT EXISTS flow_tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    flow_id BIGINT NOT NULL COMMENT '流程ID',
    task_id BIGINT NOT NULL COMMENT '任务ID',
    sequence INT NOT NULL COMMENT '执行顺序',
    is_optional TINYINT(1) DEFAULT 0 COMMENT '是否可选（可跳过）',
    allow_rollback TINYINT(1) DEFAULT 1 COMMENT '是否允许打回',
    condition_config JSON COMMENT '执行条件配置',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (flow_id) REFERENCES flows(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    UNIQUE KEY uk_flow_task_seq (flow_id, sequence),
    INDEX idx_flow_id (flow_id),
    INDEX idx_task_id (task_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='流程任务关联表';

-- 4. 作业实例表
CREATE TABLE IF NOT EXISTS jobs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    flow_id BIGINT NOT NULL COMMENT '流程ID',
    job_name VARCHAR(100) NOT NULL COMMENT '作业名称',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending/running/completed/failed/cancelled',
    current_task_seq INT COMMENT '当前执行任务序号',
    started_at TIMESTAMP NULL COMMENT '开始时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    created_by BIGINT COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (flow_id) REFERENCES flows(id),
    INDEX idx_flow_id (flow_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='作业实例表';

-- 5. 作业任务执行记录表
CREATE TABLE IF NOT EXISTS job_tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '作业ID',
    flow_task_id BIGINT NOT NULL COMMENT '流程任务ID',
    task_id BIGINT NOT NULL COMMENT '任务ID',
    sequence INT NOT NULL COMMENT '执行顺序',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending/running/completed/failed/skipped/rolled_back',
    is_skipped TINYINT(1) DEFAULT 0 COMMENT '是否跳过',
    executor_id BIGINT COMMENT '执行人ID',
    result JSON COMMENT '执行结果',
    error_message TEXT COMMENT '错误信息',
    started_at TIMESTAMP NULL COMMENT '开始时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
    FOREIGN KEY (flow_task_id) REFERENCES flow_tasks(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    INDEX idx_job_id (job_id),
    INDEX idx_status (status),
    INDEX idx_sequence (job_id, sequence)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='作业任务执行记录表';

-- 6. 作业任务日志表
CREATE TABLE IF NOT EXISTS job_task_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_task_id BIGINT NOT NULL COMMENT '作业任务ID',
    action VARCHAR(50) NOT NULL COMMENT '操作：start/complete/skip/rollback/fail',
    operator_id BIGINT COMMENT '操作人ID',
    message TEXT COMMENT '日志信息',
    metadata JSON COMMENT '元数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (job_task_id) REFERENCES job_tasks(id) ON DELETE CASCADE,
    INDEX idx_job_task_id (job_task_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='作业任务日志表';
