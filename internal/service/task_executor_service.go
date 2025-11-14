package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cfrs2005/GoWorkFlow/internal/engine"
	"github.com/cfrs2005/GoWorkFlow/internal/executor"
	"github.com/cfrs2005/GoWorkFlow/internal/models"
	"github.com/cfrs2005/GoWorkFlow/internal/repository"
	"github.com/cfrs2005/GoWorkFlow/pkg/logger"
)

// TaskExecutorService 任务执行服务
type TaskExecutorService struct {
	jobRepo        repository.JobRepository
	jobTaskRepo    repository.JobTaskRepository
	jobContextRepo repository.JobContextRepository
	taskRepo       repository.TaskRepository
	engine         engine.WorkflowEngine
}

// NewTaskExecutorService 创建任务执行服务
func NewTaskExecutorService(
	jobRepo repository.JobRepository,
	jobTaskRepo repository.JobTaskRepository,
	jobContextRepo repository.JobContextRepository,
	taskRepo repository.TaskRepository,
	workflowEngine engine.WorkflowEngine,
) *TaskExecutorService {
	return &TaskExecutorService{
		jobRepo:        jobRepo,
		jobTaskRepo:    jobTaskRepo,
		jobContextRepo: jobContextRepo,
		taskRepo:       taskRepo,
		engine:         workflowEngine,
	}
}

// ExecuteTask 自动执行任务
func (s *TaskExecutorService) ExecuteTask(ctx context.Context, jobTaskID int64) error {
	// 获取任务信息
	jobTask, err := s.jobTaskRepo.GetByID(jobTaskID)
	if err != nil {
		return fmt.Errorf("failed to get job task: %w", err)
	}

	// 获取任务定义
	task, err := s.taskRepo.GetByID(jobTask.TaskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// 只自动执行 automated 类型的任务
	if task.TaskType != models.TaskTypeAutomated {
		logger.Infof("Task %d is not automated (type: %s), skipping auto execution", task.ID, task.TaskType)
		return nil
	}

	// 标记任务开始
	if err := s.engine.StartTask(jobTaskID, 0); err != nil {
		return fmt.Errorf("failed to start task: %w", err)
	}

	logger.Infof("Starting automated execution for task %d (job task %d)", task.ID, jobTaskID)

	// 获取 Job Context
	jobContext, err := s.jobContextRepo.GetByJobID(jobTask.JobID)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to get job context: %v", err))
		jobContext = make(map[string]string)
	}

	// 从 Job Context 构建输入参数
	input := make(map[string]interface{})
	for k, v := range jobContext {
		input[k] = v
	}

	// 合并任务配置到输入
	if task.Config != nil {
		for k, v := range task.Config {
			if _, exists := input[k]; !exists {
				input[k] = v
			}
		}
	}

	// 获取执行器名称
	executorName, ok := task.Config["executor"].(string)
	if !ok || executorName == "" {
		return fmt.Errorf("task config missing 'executor' field")
	}

	// 获取执行器
	exec, err := executor.GetExecutor(executorName)
	if err != nil {
		s.engine.FailTask(jobTaskID, fmt.Sprintf("Executor not found: %s", executorName))
		return fmt.Errorf("failed to get executor: %w", err)
	}

	// 执行任务
	logger.Infof("Executing task with executor: %s", executorName)
	result, err := exec.Execute(ctx, input, jobContext)
	if err != nil {
		logger.Errorf("Task execution failed: %v", err)
		s.engine.FailTask(jobTaskID, err.Error())
		return fmt.Errorf("execution failed: %w", err)
	}

	// 保存结果到 Job Context
	if err := s.saveResultToContext(jobTask.JobID, result); err != nil {
		logger.Info(fmt.Sprintf("Failed to save result to context: %v", err))
	}

	// 标记任务完成
	taskResult := models.TaskResult(result)
	if err := s.engine.CompleteTask(jobTaskID, taskResult); err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}

	logger.Infof("Task %d (job task %d) completed successfully", task.ID, jobTaskID)

	return nil
}

// saveResultToContext 将执行结果保存到 Job Context
func (s *TaskExecutorService) saveResultToContext(jobID int64, result map[string]interface{}) error {
	for key, value := range result {
		// 将值转换为字符串
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case int, int64, float64, bool:
			valueStr = fmt.Sprintf("%v", v)
		default:
			// 复杂类型转 JSON
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				logger.Info(fmt.Sprintf("Failed to marshal value for key %s: %v", key, err))
				continue
			}
			valueStr = string(jsonBytes)
		}

		if err := s.jobContextRepo.Set(jobID, key, valueStr); err != nil {
			return fmt.Errorf("failed to set context %s: %w", key, err)
		}
	}
	return nil
}

// AutoExecuteJobTasks 自动执行作业的所有任务
func (s *TaskExecutorService) AutoExecuteJobTasks(ctx context.Context, jobID int64) error {
	logger.Infof("Starting auto execution for job %d", jobID)

	// 启动作业
	if err := s.engine.StartJob(jobID); err != nil {
		return fmt.Errorf("failed to start job: %w", err)
	}

	// 循环执行任务
	for {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 获取下一个待执行的任务
		nextTask, err := s.engine.GetNextTask(jobID)
		if err != nil {
			return fmt.Errorf("failed to get next task: %w", err)
		}

		// 如果没有下一个任务，说明流程已完成
		if nextTask == nil {
			logger.Infof("No more tasks for job %d, job completed", jobID)
			break
		}

		logger.Infof("Executing next task: job_task_id=%d, task_id=%d, sequence=%d",
			nextTask.ID, nextTask.TaskID, nextTask.Sequence)

		// 执行任务
		if err := s.ExecuteTask(ctx, nextTask.ID); err != nil {
			logger.Errorf("Failed to execute task %d: %v", nextTask.ID, err)
			// 继续尝试下一个任务，或者根据策略决定是否中断
			// 这里选择中断
			return fmt.Errorf("task execution failed: %w", err)
		}

		// 短暂等待，避免过快执行
		time.Sleep(1 * time.Second)
	}

	logger.Infof("Job %d auto execution completed", jobID)
	return nil
}
