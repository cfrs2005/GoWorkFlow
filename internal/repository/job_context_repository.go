package repository

import (
	"database/sql"
)

// JobContextRepository 作业上下文仓储接口
type JobContextRepository interface {
	GetByJobID(jobID int64) (map[string]string, error)
	Set(jobID int64, key, value string) error
	Get(jobID int64, key string) (string, error)
	Delete(jobID int64, key string) error
	DeleteByJobID(jobID int64) error
}

type jobContextRepository struct {
	db *sql.DB
}

// NewJobContextRepository 创建作业上下文仓储
func NewJobContextRepository(db *sql.DB) JobContextRepository {
	return &jobContextRepository{db: db}
}

// GetByJobID 获取作业的所有上下文数据
func (r *jobContextRepository) GetByJobID(jobID int64) (map[string]string, error) {
	query := `
		SELECT context_key, context_value
		FROM job_context
		WHERE job_id = ?
	`

	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	context := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		context[key] = value
	}

	return context, rows.Err()
}

// Set 设置上下文数据
func (r *jobContextRepository) Set(jobID int64, key, value string) error {
	query := `
		INSERT INTO job_context (job_id, context_key, context_value)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE context_value = VALUES(context_value)
	`

	_, err := r.db.Exec(query, jobID, key, value)
	return err
}

// Get 获取单个上下文值
func (r *jobContextRepository) Get(jobID int64, key string) (string, error) {
	query := `
		SELECT context_value
		FROM job_context
		WHERE job_id = ? AND context_key = ?
	`

	var value string
	err := r.db.QueryRow(query, jobID, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}

	return value, err
}

// Delete 删除单个上下文键
func (r *jobContextRepository) Delete(jobID int64, key string) error {
	query := `DELETE FROM job_context WHERE job_id = ? AND context_key = ?`
	_, err := r.db.Exec(query, jobID, key)
	return err
}

// DeleteByJobID 删除作业的所有上下文数据
func (r *jobContextRepository) DeleteByJobID(jobID int64) error {
	query := `DELETE FROM job_context WHERE job_id = ?`
	_, err := r.db.Exec(query, jobID)
	return err
}
