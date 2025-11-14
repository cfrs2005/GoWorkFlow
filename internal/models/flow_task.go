package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ConditionConfig 执行条件配置
type ConditionConfig map[string]interface{}

// Value 实现 driver.Valuer 接口
func (cc ConditionConfig) Value() (driver.Value, error) {
	if cc == nil {
		return nil, nil
	}
	return json.Marshal(cc)
}

// Scan 实现 sql.Scanner 接口
func (cc *ConditionConfig) Scan(value interface{}) error {
	if value == nil {
		*cc = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, cc)
}

// FlowTask 流程任务关联模型
type FlowTask struct {
	ID              int64           `json:"id"`
	FlowID          int64           `json:"flow_id"`
	TaskID          int64           `json:"task_id"`
	Sequence        int             `json:"sequence"`
	IsOptional      bool            `json:"is_optional"`
	AllowRollback   bool            `json:"allow_rollback"`
	ConditionConfig ConditionConfig `json:"condition_config"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// 关联数据
	Task *Task `json:"task,omitempty"`
}

// TableName 返回表名
func (FlowTask) TableName() string {
	return "flow_tasks"
}
