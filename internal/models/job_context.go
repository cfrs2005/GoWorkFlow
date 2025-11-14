package models

import (
	"time"
)

// JobContext 作业上下文模型
type JobContext struct {
	ID           int64     `json:"id"`
	JobID        int64     `json:"job_id"`
	ContextKey   string    `json:"context_key"`
	ContextValue string    `json:"context_value"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 返回表名
func (JobContext) TableName() string {
	return "job_context"
}
