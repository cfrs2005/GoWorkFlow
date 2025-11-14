package service

import (
	"database/sql"
	"fmt"

	"github.com/cfrs2005/GoWorkFlow/internal/engine"
	"github.com/cfrs2005/GoWorkFlow/internal/models"
	"github.com/cfrs2005/GoWorkFlow/internal/repository"
)

// WorkflowService 工作流服务接口
type WorkflowService interface {
	// Task 管理
	CreateTask(task *models.Task) error
	GetTask(id int64) (*models.Task, error)
	ListTasks(limit, offset int) ([]models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id int64) error

	// Flow 管理
	CreateFlow(flow *models.Flow, taskIDs []int64) error
	GetFlow(id int64) (*models.Flow, error)
	GetFlowWithTasks(id int64) (*models.Flow, []models.FlowTask, error)
	ListFlows(limit, offset int) ([]models.Flow, error)
	UpdateFlow(flow *models.Flow) error
	DeleteFlow(id int64) error
	AddTaskToFlow(flowID, taskID int64, sequence int, isOptional, allowRollback bool) error

	// Job 管理
	CreateJob(flowID int64, jobName string, createdBy int64) (*models.Job, error)
	GetJob(id int64) (*models.Job, error)
	GetJobWithTasks(id int64) (*models.Job, []models.JobTask, error)
	ListJobs(limit, offset int) ([]models.Job, error)
	StartJob(jobID int64) error

	// JobTask 操作
	StartTask(jobTaskID int64, executorID int64) error
	CompleteTask(jobTaskID int64, result models.TaskResult) error
	FailTask(jobTaskID int64, errorMessage string) error
	SkipTask(jobTaskID int64, operatorID int64) error
	RollbackTask(jobTaskID int64, operatorID int64, targetSequence int) error
	GetNextTask(jobID int64) (*models.JobTask, error)
}

type workflowService struct {
	db           *sql.DB
	taskRepo     repository.TaskRepository
	flowRepo     repository.FlowRepository
	flowTaskRepo repository.FlowTaskRepository
	jobRepo      repository.JobRepository
	jobTaskRepo  repository.JobTaskRepository
	engine       engine.WorkflowEngine
}

// NewWorkflowService 创建工作流服务
func NewWorkflowService(
	db *sql.DB,
	taskRepo repository.TaskRepository,
	flowRepo repository.FlowRepository,
	flowTaskRepo repository.FlowTaskRepository,
	jobRepo repository.JobRepository,
	jobTaskRepo repository.JobTaskRepository,
	engine engine.WorkflowEngine,
) WorkflowService {
	return &workflowService{
		db:           db,
		taskRepo:     taskRepo,
		flowRepo:     flowRepo,
		flowTaskRepo: flowTaskRepo,
		jobRepo:      jobRepo,
		jobTaskRepo:  jobTaskRepo,
		engine:       engine,
	}
}

// Task 管理方法

func (s *workflowService) CreateTask(task *models.Task) error {
	return s.taskRepo.Create(task)
}

func (s *workflowService) GetTask(id int64) (*models.Task, error) {
	return s.taskRepo.GetByID(id)
}

func (s *workflowService) ListTasks(limit, offset int) ([]models.Task, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.taskRepo.List(limit, offset)
}

func (s *workflowService) UpdateTask(task *models.Task) error {
	return s.taskRepo.Update(task)
}

func (s *workflowService) DeleteTask(id int64) error {
	return s.taskRepo.Delete(id)
}

// Flow 管理方法

func (s *workflowService) CreateFlow(flow *models.Flow, taskIDs []int64) error {
	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 创建流程
	if err := s.flowRepo.Create(flow); err != nil {
		return err
	}

	// 添加任务到流程
	for i, taskID := range taskIDs {
		flowTask := &models.FlowTask{
			FlowID:        flow.ID,
			TaskID:        taskID,
			Sequence:      i + 1,
			IsOptional:    false,
			AllowRollback: true,
		}
		if err := s.flowTaskRepo.Create(flowTask); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *workflowService) GetFlow(id int64) (*models.Flow, error) {
	return s.flowRepo.GetByID(id)
}

func (s *workflowService) GetFlowWithTasks(id int64) (*models.Flow, []models.FlowTask, error) {
	return s.flowRepo.GetFlowWithTasks(id)
}

func (s *workflowService) ListFlows(limit, offset int) ([]models.Flow, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.flowRepo.List(limit, offset)
}

func (s *workflowService) UpdateFlow(flow *models.Flow) error {
	return s.flowRepo.Update(flow)
}

func (s *workflowService) DeleteFlow(id int64) error {
	return s.flowRepo.Delete(id)
}

func (s *workflowService) AddTaskToFlow(flowID, taskID int64, sequence int, isOptional, allowRollback bool) error {
	flowTask := &models.FlowTask{
		FlowID:        flowID,
		TaskID:        taskID,
		Sequence:      sequence,
		IsOptional:    isOptional,
		AllowRollback: allowRollback,
	}
	return s.flowTaskRepo.Create(flowTask)
}

// Job 管理方法

func (s *workflowService) CreateJob(flowID int64, jobName string, createdBy int64) (*models.Job, error) {
	return s.engine.CreateJob(flowID, jobName, createdBy)
}

func (s *workflowService) GetJob(id int64) (*models.Job, error) {
	return s.jobRepo.GetByID(id)
}

func (s *workflowService) GetJobWithTasks(id int64) (*models.Job, []models.JobTask, error) {
	return s.jobRepo.GetJobWithTasks(id)
}

func (s *workflowService) ListJobs(limit, offset int) ([]models.Job, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.jobRepo.List(limit, offset)
}

func (s *workflowService) StartJob(jobID int64) error {
	return s.engine.StartJob(jobID)
}

// JobTask 操作方法

func (s *workflowService) StartTask(jobTaskID int64, executorID int64) error {
	return s.engine.StartTask(jobTaskID, executorID)
}

func (s *workflowService) CompleteTask(jobTaskID int64, result models.TaskResult) error {
	return s.engine.CompleteTask(jobTaskID, result)
}

func (s *workflowService) FailTask(jobTaskID int64, errorMessage string) error {
	return s.engine.FailTask(jobTaskID, errorMessage)
}

func (s *workflowService) SkipTask(jobTaskID int64, operatorID int64) error {
	return s.engine.SkipTask(jobTaskID, operatorID)
}

func (s *workflowService) RollbackTask(jobTaskID int64, operatorID int64, targetSequence int) error {
	return s.engine.RollbackTask(jobTaskID, operatorID, targetSequence)
}

func (s *workflowService) GetNextTask(jobID int64) (*models.JobTask, error) {
	return s.engine.GetNextTask(jobID)
}
