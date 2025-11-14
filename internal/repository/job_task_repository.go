package repository

import (
	"database/sql"
	"fmt"

	"github.com/cfrs2005/GoWorkFlow/internal/models"
)

// JobTaskRepository 作业任务仓储接口
type JobTaskRepository interface {
	Create(jobTask *models.JobTask) error
	GetByID(id int64) (*models.JobTask, error)
	GetByJobID(jobID int64) ([]models.JobTask, error)
	GetBySequence(jobID int64, sequence int) (*models.JobTask, error)
	Update(jobTask *models.JobTask) error
	UpdateStatus(id int64, status models.JobTaskStatus) error
	BatchCreate(jobTasks []models.JobTask) error
}

type jobTaskRepository struct {
	db *sql.DB
}

// NewJobTaskRepository 创建作业任务仓储
func NewJobTaskRepository(db *sql.DB) JobTaskRepository {
	return &jobTaskRepository{db: db}
}

// Create 创建作业任务
func (r *jobTaskRepository) Create(jobTask *models.JobTask) error {
	query := `
		INSERT INTO job_tasks (job_id, flow_task_id, task_id, sequence, status, is_skipped)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		jobTask.JobID, jobTask.FlowTaskID, jobTask.TaskID,
		jobTask.Sequence, jobTask.Status, jobTask.IsSkipped,
	)
	if err != nil {
		return fmt.Errorf("failed to create job task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	jobTask.ID = id
	return nil
}

// GetByID 根据ID获取作业任务
func (r *jobTaskRepository) GetByID(id int64) (*models.JobTask, error) {
	query := `
		SELECT id, job_id, flow_task_id, task_id, sequence, status, is_skipped,
		       executor_id, result, error_message, started_at, completed_at,
		       created_at, updated_at
		FROM job_tasks
		WHERE id = ?
	`
	jobTask := &models.JobTask{}
	err := r.db.QueryRow(query, id).Scan(
		&jobTask.ID, &jobTask.JobID, &jobTask.FlowTaskID, &jobTask.TaskID,
		&jobTask.Sequence, &jobTask.Status, &jobTask.IsSkipped, &jobTask.ExecutorID,
		&jobTask.Result, &jobTask.ErrorMessage, &jobTask.StartedAt, &jobTask.CompletedAt,
		&jobTask.CreatedAt, &jobTask.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job task not found")
		}
		return nil, fmt.Errorf("failed to get job task: %w", err)
	}

	return jobTask, nil
}

// GetByJobID 根据作业ID获取所有作业任务
func (r *jobTaskRepository) GetByJobID(jobID int64) ([]models.JobTask, error) {
	query := `
		SELECT id, job_id, flow_task_id, task_id, sequence, status, is_skipped,
		       executor_id, result, error_message, started_at, completed_at,
		       created_at, updated_at
		FROM job_tasks
		WHERE job_id = ?
		ORDER BY sequence ASC
	`
	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job tasks: %w", err)
	}
	defer rows.Close()

	var jobTasks []models.JobTask
	for rows.Next() {
		var jobTask models.JobTask
		if err := rows.Scan(
			&jobTask.ID, &jobTask.JobID, &jobTask.FlowTaskID, &jobTask.TaskID,
			&jobTask.Sequence, &jobTask.Status, &jobTask.IsSkipped, &jobTask.ExecutorID,
			&jobTask.Result, &jobTask.ErrorMessage, &jobTask.StartedAt, &jobTask.CompletedAt,
			&jobTask.CreatedAt, &jobTask.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan job task: %w", err)
		}
		jobTasks = append(jobTasks, jobTask)
	}

	return jobTasks, nil
}

// GetBySequence 根据作业ID和序号获取作业任务
func (r *jobTaskRepository) GetBySequence(jobID int64, sequence int) (*models.JobTask, error) {
	query := `
		SELECT id, job_id, flow_task_id, task_id, sequence, status, is_skipped,
		       executor_id, result, error_message, started_at, completed_at,
		       created_at, updated_at
		FROM job_tasks
		WHERE job_id = ? AND sequence = ?
	`
	jobTask := &models.JobTask{}
	err := r.db.QueryRow(query, jobID, sequence).Scan(
		&jobTask.ID, &jobTask.JobID, &jobTask.FlowTaskID, &jobTask.TaskID,
		&jobTask.Sequence, &jobTask.Status, &jobTask.IsSkipped, &jobTask.ExecutorID,
		&jobTask.Result, &jobTask.ErrorMessage, &jobTask.StartedAt, &jobTask.CompletedAt,
		&jobTask.CreatedAt, &jobTask.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job task not found")
		}
		return nil, fmt.Errorf("failed to get job task: %w", err)
	}

	return jobTask, nil
}

// Update 更新作业任务
func (r *jobTaskRepository) Update(jobTask *models.JobTask) error {
	query := `
		UPDATE job_tasks
		SET status = ?, is_skipped = ?, executor_id = ?, result = ?,
		    error_message = ?, started_at = ?, completed_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		jobTask.Status, jobTask.IsSkipped, jobTask.ExecutorID, jobTask.Result,
		jobTask.ErrorMessage, jobTask.StartedAt, jobTask.CompletedAt, jobTask.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update job task: %w", err)
	}

	return nil
}

// UpdateStatus 更新作业任务状态
func (r *jobTaskRepository) UpdateStatus(id int64, status models.JobTaskStatus) error {
	query := `UPDATE job_tasks SET status = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update job task status: %w", err)
	}

	return nil
}

// BatchCreate 批量创建作业任务
func (r *jobTaskRepository) BatchCreate(jobTasks []models.JobTask) error {
	if len(jobTasks) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO job_tasks (job_id, flow_task_id, task_id, sequence, status, is_skipped)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := range jobTasks {
		result, err := stmt.Exec(
			jobTasks[i].JobID, jobTasks[i].FlowTaskID, jobTasks[i].TaskID,
			jobTasks[i].Sequence, jobTasks[i].Status, jobTasks[i].IsSkipped,
		)
		if err != nil {
			return fmt.Errorf("failed to insert job task: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		jobTasks[i].ID = id
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
