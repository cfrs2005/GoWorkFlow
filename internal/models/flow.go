package models

import (
	"time"
)

// Flow 流程定义模型
type Flow struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 返回表名
func (Flow) TableName() string {
	return "flows"
}
