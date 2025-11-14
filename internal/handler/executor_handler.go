package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cfrs2005/GoWorkFlow/internal/service"
	"github.com/cfrs2005/GoWorkFlow/pkg/logger"
	"github.com/cfrs2005/GoWorkFlow/pkg/response"
)

// ExecutorHandler 任务执行处理器
type ExecutorHandler struct {
	taskExecutorService *service.TaskExecutorService
}

// NewExecutorHandler 创建任务执行处理器
func NewExecutorHandler(taskExecutorService *service.TaskExecutorService) *ExecutorHandler {
	return &ExecutorHandler{
		taskExecutorService: taskExecutorService,
	}
}

// AutoExecuteJob 自动执行作业的所有任务
// POST /api/jobs/auto-execute
func (h *ExecutorHandler) AutoExecuteJob(w http.ResponseWriter, r *http.Request) {
	var req struct {
		JobID int64 `json:"job_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.JobID == 0 {
		response.Error(w, http.StatusBadRequest, "job_id is required")
		return
	}

	logger.Infof("Starting auto execution for job %d", req.JobID)

	// 创建超时上下文（最多30分钟）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// 在后台执行
	go func() {
		if err := h.taskExecutorService.AutoExecuteJobTasks(ctx, req.JobID); err != nil {
			logger.Errorf("Auto execution failed for job %d: %v", req.JobID, err)
		}
	}()

	response.Success(w, map[string]interface{}{
		"message": "Auto execution started",
		"job_id":  req.JobID,
	})
}

// ExecuteTask 执行单个任务
// POST /api/tasks/execute
func (h *ExecutorHandler) ExecuteTask(w http.ResponseWriter, r *http.Request) {
	jobTaskIDStr := r.URL.Query().Get("job_task_id")
	if jobTaskIDStr == "" {
		response.Error(w, http.StatusBadRequest, "job_task_id is required")
		return
	}

	jobTaskID, err := strconv.ParseInt(jobTaskIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid job_task_id")
		return
	}

	// 创建超时上下文（最多10分钟）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// 在后台执行
	go func() {
		if err := h.taskExecutorService.ExecuteTask(ctx, jobTaskID); err != nil {
			logger.Errorf("Task execution failed for job_task %d: %v", jobTaskID, err)
		}
	}()

	response.Success(w, map[string]interface{}{
		"message":     "Task execution started",
		"job_task_id": jobTaskID,
	})
}
