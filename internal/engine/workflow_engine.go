package engine

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cfrs2005/GoWorkFlow/internal/models"
	"github.com/cfrs2005/GoWorkFlow/internal/repository"
)

// WorkflowEngine 工作流引擎接口
type WorkflowEngine interface {
	// CreateJob 创建作业实例
	CreateJob(flowID int64, jobName string, createdBy int64) (*models.Job, error)

	// StartJob 启动作业
	StartJob(jobID int64) error

	// StartTask 开始执行任务
	StartTask(jobTaskID int64, executorID int64) error

	// CompleteTask 完成任务
	CompleteTask(jobTaskID int64, result models.TaskResult) error

	// FailTask 任务失败
	FailTask(jobTaskID int64, errorMessage string) error

	// SkipTask 跳过任务
	SkipTask(jobTaskID int64, operatorID int64) error

	// RollbackTask 打回任务
	RollbackTask(jobTaskID int64, operatorID int64, targetSequence int) error

	// GetNextTask 获取下一个待执行的任务
	GetNextTask(jobID int64) (*models.JobTask, error)

	// GetCurrentTask 获取当前执行中的任务
	GetCurrentTask(jobID int64) (*models.JobTask, error)
}

type workflowEngine struct {
	db              *sql.DB
	jobRepo         repository.JobRepository
	jobTaskRepo     repository.JobTaskRepository
	flowRepo        repository.FlowRepository
	flowTaskRepo    repository.FlowTaskRepository
}

// NewWorkflowEngine 创建工作流引擎
func NewWorkflowEngine(
	db *sql.DB,
	jobRepo repository.JobRepository,
	jobTaskRepo repository.JobTaskRepository,
	flowRepo repository.FlowRepository,
	flowTaskRepo repository.FlowTaskRepository,
) WorkflowEngine {
	return &workflowEngine{
		db:           db,
		jobRepo:      jobRepo,
		jobTaskRepo:  jobTaskRepo,
		flowRepo:     flowRepo,
		flowTaskRepo: flowTaskRepo,
	}
}

// CreateJob 创建作业实例
func (e *workflowEngine) CreateJob(flowID int64, jobName string, createdBy int64) (*models.Job, error) {
	// 获取流程及其任务
	flow, flowTasks, err := e.flowRepo.GetFlowWithTasks(flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get flow: %w", err)
	}

	if !flow.IsActive {
		return nil, fmt.Errorf("flow is not active")
	}

	if len(flowTasks) == 0 {
		return nil, fmt.Errorf("flow has no tasks")
	}

	// 开始事务
	tx, err := e.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 创建作业
	job := &models.Job{
		FlowID:    flowID,
		JobName:   jobName,
		Status:    models.JobStatusPending,
		CreatedBy: createdBy,
	}

	if err := e.jobRepo.Create(job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	// 为每个流程任务创建作业任务
	jobTasks := make([]models.JobTask, len(flowTasks))
	for i, ft := range flowTasks {
		jobTasks[i] = models.JobTask{
			JobID:      job.ID,
			FlowTaskID: ft.ID,
			TaskID:     ft.TaskID,
			Sequence:   ft.Sequence,
			Status:     models.JobTaskStatusPending,
			IsSkipped:  false,
		}
	}

	if err := e.jobTaskRepo.BatchCreate(jobTasks); err != nil {
		return nil, fmt.Errorf("failed to create job tasks: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return job, nil
}

// StartJob 启动作业
func (e *workflowEngine) StartJob(jobID int64) error {
	job, err := e.jobRepo.GetByID(jobID)
	if err != nil {
		return err
	}

	if job.Status != models.JobStatusPending {
		return fmt.Errorf("job is not in pending status")
	}

	// 更新作业状态
	job.Status = models.JobStatusRunning
	job.StartedAt = sql.NullTime{Time: time.Now(), Valid: true}
	job.CurrentTaskSeq = sql.NullInt64{Int64: 1, Valid: true}

	return e.jobRepo.Update(job)
}

// StartTask 开始执行任务
func (e *workflowEngine) StartTask(jobTaskID int64, executorID int64) error {
	jobTask, err := e.jobTaskRepo.GetByID(jobTaskID)
	if err != nil {
		return err
	}

	if jobTask.Status != models.JobTaskStatusPending {
		return fmt.Errorf("job task is not in pending status")
	}

	// 更新任务状态
	jobTask.Status = models.JobTaskStatusRunning
	jobTask.StartedAt = sql.NullTime{Time: time.Now(), Valid: true}
	jobTask.ExecutorID = sql.NullInt64{Int64: executorID, Valid: true}

	return e.jobTaskRepo.Update(jobTask)
}

// CompleteTask 完成任务
func (e *workflowEngine) CompleteTask(jobTaskID int64, result models.TaskResult) error {
	jobTask, err := e.jobTaskRepo.GetByID(jobTaskID)
	if err != nil {
		return err
	}

	if jobTask.Status != models.JobTaskStatusRunning {
		return fmt.Errorf("job task is not in running status")
	}

	// 更新任务状态
	jobTask.Status = models.JobTaskStatusCompleted
	jobTask.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}
	jobTask.Result = result

	if err := e.jobTaskRepo.Update(jobTask); err != nil {
		return err
	}

	// 检查是否有下一个任务
	nextTask, err := e.GetNextTask(jobTask.JobID)
	if err != nil || nextTask == nil {
		// 没有下一个任务，完成作业
		return e.completeJob(jobTask.JobID)
	}

	// 更新作业的当前任务序号
	job, err := e.jobRepo.GetByID(jobTask.JobID)
	if err != nil {
		return err
	}

	job.CurrentTaskSeq = sql.NullInt64{Int64: int64(nextTask.Sequence), Valid: true}
	return e.jobRepo.Update(job)
}

// FailTask 任务失败
func (e *workflowEngine) FailTask(jobTaskID int64, errorMessage string) error {
	jobTask, err := e.jobTaskRepo.GetByID(jobTaskID)
	if err != nil {
		return err
	}

	if jobTask.Status != models.JobTaskStatusRunning {
		return fmt.Errorf("job task is not in running status")
	}

	// 更新任务状态
	jobTask.Status = models.JobTaskStatusFailed
	jobTask.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}
	jobTask.ErrorMessage = errorMessage

	if err := e.jobTaskRepo.Update(jobTask); err != nil {
		return err
	}

	// 更新作业状态为失败
	return e.jobRepo.UpdateStatus(jobTask.JobID, models.JobStatusFailed)
}

// SkipTask 跳过任务
func (e *workflowEngine) SkipTask(jobTaskID int64, operatorID int64) error {
	jobTask, err := e.jobTaskRepo.GetByID(jobTaskID)
	if err != nil {
		return err
	}

	if jobTask.Status != models.JobTaskStatusPending {
		return fmt.Errorf("job task is not in pending status")
	}

	// 检查任务是否可跳过
	flowTask, err := e.flowTaskRepo.GetByID(jobTask.FlowTaskID)
	if err != nil {
		return err
	}

	if !flowTask.IsOptional {
		return fmt.Errorf("task is not optional and cannot be skipped")
	}

	// 更新任务状态
	jobTask.Status = models.JobTaskStatusSkipped
	jobTask.IsSkipped = true
	jobTask.ExecutorID = sql.NullInt64{Int64: operatorID, Valid: true}
	jobTask.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := e.jobTaskRepo.Update(jobTask); err != nil {
		return err
	}

	// 检查是否有下一个任务
	nextTask, err := e.GetNextTask(jobTask.JobID)
	if err != nil || nextTask == nil {
		// 没有下一个任务，完成作业
		return e.completeJob(jobTask.JobID)
	}

	// 更新作业的当前任务序号
	job, err := e.jobRepo.GetByID(jobTask.JobID)
	if err != nil {
		return err
	}

	job.CurrentTaskSeq = sql.NullInt64{Int64: int64(nextTask.Sequence), Valid: true}
	return e.jobRepo.Update(job)
}

// RollbackTask 打回任务
func (e *workflowEngine) RollbackTask(jobTaskID int64, operatorID int64, targetSequence int) error {
	jobTask, err := e.jobTaskRepo.GetByID(jobTaskID)
	if err != nil {
		return err
	}

	// 检查任务是否允许打回
	flowTask, err := e.flowTaskRepo.GetByID(jobTask.FlowTaskID)
	if err != nil {
		return err
	}

	if !flowTask.AllowRollback {
		return fmt.Errorf("task does not allow rollback")
	}

	// 获取目标任务
	targetTask, err := e.jobTaskRepo.GetBySequence(jobTask.JobID, targetSequence)
	if err != nil {
		return fmt.Errorf("target task not found: %w", err)
	}

	if targetSequence >= jobTask.Sequence {
		return fmt.Errorf("can only rollback to previous tasks")
	}

	// 更新当前任务状态为已打回
	jobTask.Status = models.JobTaskStatusRolledBack
	jobTask.ExecutorID = sql.NullInt64{Int64: operatorID, Valid: true}

	if err := e.jobTaskRepo.Update(jobTask); err != nil {
		return err
	}

	// 重置目标任务及之后的所有任务状态
	jobTasks, err := e.jobTaskRepo.GetByJobID(jobTask.JobID)
	if err != nil {
		return err
	}

	for i := range jobTasks {
		if jobTasks[i].Sequence >= targetSequence {
			jobTasks[i].Status = models.JobTaskStatusPending
			jobTasks[i].IsSkipped = false
			jobTasks[i].ExecutorID = sql.NullInt64{}
			jobTasks[i].Result = nil
			jobTasks[i].ErrorMessage = ""
			jobTasks[i].StartedAt = sql.NullTime{}
			jobTasks[i].CompletedAt = sql.NullTime{}

			if err := e.jobTaskRepo.Update(&jobTasks[i]); err != nil {
				return err
			}
		}
	}

	// 更新作业的当前任务序号
	job, err := e.jobRepo.GetByID(jobTask.JobID)
	if err != nil {
		return err
	}

	job.CurrentTaskSeq = sql.NullInt64{Int64: int64(targetTask.Sequence), Valid: true}
	job.Status = models.JobStatusRunning

	return e.jobRepo.Update(job)
}

// GetNextTask 获取下一个待执行的任务
func (e *workflowEngine) GetNextTask(jobID int64) (*models.JobTask, error) {
	jobTasks, err := e.jobTaskRepo.GetByJobID(jobID)
	if err != nil {
		return nil, err
	}

	for i := range jobTasks {
		if jobTasks[i].Status == models.JobTaskStatusPending {
			return &jobTasks[i], nil
		}
	}

	return nil, nil
}

// GetCurrentTask 获取当前执行中的任务
func (e *workflowEngine) GetCurrentTask(jobID int64) (*models.JobTask, error) {
	job, err := e.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, err
	}

	if !job.CurrentTaskSeq.Valid {
		return nil, fmt.Errorf("no current task")
	}

	return e.jobTaskRepo.GetBySequence(jobID, int(job.CurrentTaskSeq.Int64))
}

// completeJob 完成作业
func (e *workflowEngine) completeJob(jobID int64) error {
	job, err := e.jobRepo.GetByID(jobID)
	if err != nil {
		return err
	}

	job.Status = models.JobStatusCompleted
	job.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return e.jobRepo.Update(job)
}
