package models

import (
	"database/sql"
	"time"
)

// JobStatus 作业状态
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"   // 待执行
	JobStatusRunning   JobStatus = "running"   // 执行中
	JobStatusCompleted JobStatus = "completed" // 已完成
	JobStatusFailed    JobStatus = "failed"    // 失败
	JobStatusCancelled JobStatus = "cancelled" // 已取消
)

// Job 作业实例模型
type Job struct {
	ID             int64        `json:"id"`
	FlowID         int64        `json:"flow_id"`
	JobName        string       `json:"job_name"`
	Status         JobStatus    `json:"status"`
	CurrentTaskSeq sql.NullInt64 `json:"current_task_seq"`
	StartedAt      sql.NullTime `json:"started_at"`
	CompletedAt    sql.NullTime `json:"completed_at"`
	CreatedBy      int64        `json:"created_by"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`

	// 关联数据
	Flow     *Flow      `json:"flow,omitempty"`
	JobTasks []JobTask  `json:"job_tasks,omitempty"`
}

// TableName 返回表名
func (Job) TableName() string {
	return "jobs"
}
