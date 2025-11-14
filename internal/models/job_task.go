package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JobTaskStatus 作业任务状态
type JobTaskStatus string

const (
	JobTaskStatusPending    JobTaskStatus = "pending"     // 待执行
	JobTaskStatusRunning    JobTaskStatus = "running"     // 执行中
	JobTaskStatusCompleted  JobTaskStatus = "completed"   // 已完成
	JobTaskStatusFailed     JobTaskStatus = "failed"      // 失败
	JobTaskStatusSkipped    JobTaskStatus = "skipped"     // 已跳过
	JobTaskStatusRolledBack JobTaskStatus = "rolled_back" // 已打回
)

// TaskResult 任务执行结果
type TaskResult map[string]interface{}

// Value 实现 driver.Valuer 接口
func (tr TaskResult) Value() (driver.Value, error) {
	if tr == nil {
		return nil, nil
	}
	return json.Marshal(tr)
}

// Scan 实现 sql.Scanner 接口
func (tr *TaskResult) Scan(value interface{}) error {
	if value == nil {
		*tr = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, tr)
}

// JobTask 作业任务执行记录模型
type JobTask struct {
	ID           int64         `json:"id"`
	JobID        int64         `json:"job_id"`
	FlowTaskID   int64         `json:"flow_task_id"`
	TaskID       int64         `json:"task_id"`
	Sequence     int           `json:"sequence"`
	Status       JobTaskStatus `json:"status"`
	IsSkipped    bool          `json:"is_skipped"`
	ExecutorID   sql.NullInt64 `json:"executor_id"`
	Result       TaskResult    `json:"result"`
	ErrorMessage string        `json:"error_message"`
	StartedAt    sql.NullTime  `json:"started_at"`
	CompletedAt  sql.NullTime  `json:"completed_at"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`

	// 关联数据
	Task     *Task     `json:"task,omitempty"`
	FlowTask *FlowTask `json:"flow_task,omitempty"`
}

// TableName 返回表名
func (JobTask) TableName() string {
	return "job_tasks"
}
