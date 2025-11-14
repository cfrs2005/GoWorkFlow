package repository

import (
	"database/sql"
	"fmt"

	"github.com/cfrs2005/GoWorkFlow/internal/models"
)

// JobRepository 作业仓储接口
type JobRepository interface {
	Create(job *models.Job) error
	GetByID(id int64) (*models.Job, error)
	List(limit, offset int) ([]models.Job, error)
	Update(job *models.Job) error
	UpdateStatus(jobID int64, status models.JobStatus) error
	GetJobWithTasks(jobID int64) (*models.Job, []models.JobTask, error)
}

type jobRepository struct {
	db *sql.DB
}

// NewJobRepository 创建作业仓储
func NewJobRepository(db *sql.DB) JobRepository {
	return &jobRepository{db: db}
}

// Create 创建作业
func (r *jobRepository) Create(job *models.Job) error {
	query := `
		INSERT INTO jobs (flow_id, job_name, status, current_task_seq, created_by)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, job.FlowID, job.JobName, job.Status, job.CurrentTaskSeq, job.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	job.ID = id
	return nil
}

// GetByID 根据ID获取作业
func (r *jobRepository) GetByID(id int64) (*models.Job, error) {
	query := `
		SELECT id, flow_id, job_name, status, current_task_seq, started_at, completed_at,
		       created_by, created_at, updated_at
		FROM jobs
		WHERE id = ?
	`
	job := &models.Job{}
	err := r.db.QueryRow(query, id).Scan(
		&job.ID, &job.FlowID, &job.JobName, &job.Status, &job.CurrentTaskSeq,
		&job.StartedAt, &job.CompletedAt, &job.CreatedBy, &job.CreatedAt, &job.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return job, nil
}

// List 获取作业列表
func (r *jobRepository) List(limit, offset int) ([]models.Job, error) {
	query := `
		SELECT id, flow_id, job_name, status, current_task_seq, started_at, completed_at,
		       created_by, created_at, updated_at
		FROM jobs
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		if err := rows.Scan(
			&job.ID, &job.FlowID, &job.JobName, &job.Status, &job.CurrentTaskSeq,
			&job.StartedAt, &job.CompletedAt, &job.CreatedBy, &job.CreatedAt, &job.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// Update 更新作业
func (r *jobRepository) Update(job *models.Job) error {
	query := `
		UPDATE jobs
		SET status = ?, current_task_seq = ?, started_at = ?, completed_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, job.Status, job.CurrentTaskSeq, job.StartedAt, job.CompletedAt, job.ID)
	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	return nil
}

// UpdateStatus 更新作业状态
func (r *jobRepository) UpdateStatus(jobID int64, status models.JobStatus) error {
	query := `UPDATE jobs SET status = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, jobID)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}

// GetJobWithTasks 获取作业及其任务
func (r *jobRepository) GetJobWithTasks(jobID int64) (*models.Job, []models.JobTask, error) {
	// 获取作业
	job, err := r.GetByID(jobID)
	if err != nil {
		return nil, nil, err
	}

	// 获取作业任务
	query := `
		SELECT jt.id, jt.job_id, jt.flow_task_id, jt.task_id, jt.sequence, jt.status,
		       jt.is_skipped, jt.executor_id, jt.result, jt.error_message,
		       jt.started_at, jt.completed_at, jt.created_at, jt.updated_at,
		       t.id, t.name, t.description, t.task_type, t.config, t.is_active,
		       t.created_at, t.updated_at
		FROM job_tasks jt
		INNER JOIN tasks t ON jt.task_id = t.id
		WHERE jt.job_id = ?
		ORDER BY jt.sequence ASC
	`
	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get job tasks: %w", err)
	}
	defer rows.Close()

	var jobTasks []models.JobTask
	for rows.Next() {
		var jobTask models.JobTask
		var task models.Task
		if err := rows.Scan(
			&jobTask.ID, &jobTask.JobID, &jobTask.FlowTaskID, &jobTask.TaskID,
			&jobTask.Sequence, &jobTask.Status, &jobTask.IsSkipped, &jobTask.ExecutorID,
			&jobTask.Result, &jobTask.ErrorMessage, &jobTask.StartedAt, &jobTask.CompletedAt,
			&jobTask.CreatedAt, &jobTask.UpdatedAt,
			&task.ID, &task.Name, &task.Description, &task.TaskType,
			&task.Config, &task.IsActive, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, nil, fmt.Errorf("failed to scan job task: %w", err)
		}
		jobTask.Task = &task
		jobTasks = append(jobTasks, jobTask)
	}

	return job, jobTasks, nil
}
