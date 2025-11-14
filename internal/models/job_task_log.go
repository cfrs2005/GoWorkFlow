package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// LogAction 日志操作类型
type LogAction string

const (
	LogActionStart    LogAction = "start"    // 开始
	LogActionComplete LogAction = "complete" // 完成
	LogActionSkip     LogAction = "skip"     // 跳过
	LogActionRollback LogAction = "rollback" // 打回
	LogActionFail     LogAction = "fail"     // 失败
)

// LogMetadata 日志元数据
type LogMetadata map[string]interface{}

// Value 实现 driver.Valuer 接口
func (lm LogMetadata) Value() (driver.Value, error) {
	if lm == nil {
		return nil, nil
	}
	return json.Marshal(lm)
}

// Scan 实现 sql.Scanner 接口
func (lm *LogMetadata) Scan(value interface{}) error {
	if value == nil {
		*lm = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, lm)
}

// JobTaskLog 作业任务日志模型
type JobTaskLog struct {
	ID         int64          `json:"id"`
	JobTaskID  int64          `json:"job_task_id"`
	Action     LogAction      `json:"action"`
	OperatorID sql.NullInt64  `json:"operator_id"`
	Message    string         `json:"message"`
	Metadata   LogMetadata    `json:"metadata"`
	CreatedAt  time.Time      `json:"created_at"`
}

// TableName 返回表名
func (JobTaskLog) TableName() string {
	return "job_task_logs"
}
