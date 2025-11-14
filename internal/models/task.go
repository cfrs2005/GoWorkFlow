package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypeManual    TaskType = "manual"    // 手动任务
	TaskTypeAutomated TaskType = "automated" // 自动化任务
	TaskTypeApproval  TaskType = "approval"  // 审批任务
)

// TaskConfig 任务配置
type TaskConfig map[string]interface{}

// Value 实现 driver.Valuer 接口
func (tc TaskConfig) Value() (driver.Value, error) {
	if tc == nil {
		return nil, nil
	}
	return json.Marshal(tc)
}

// Scan 实现 sql.Scanner 接口
func (tc *TaskConfig) Scan(value interface{}) error {
	if value == nil {
		*tc = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, tc)
}

// Task 任务定义模型
type Task struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	TaskType    TaskType   `json:"task_type"`
	Config      TaskConfig `json:"config"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName 返回表名
func (Task) TableName() string {
	return "tasks"
}
