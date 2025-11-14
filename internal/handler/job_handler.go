package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cfrs2005/GoWorkFlow/internal/models"
	"github.com/cfrs2005/GoWorkFlow/internal/service"
	"github.com/cfrs2005/GoWorkFlow/pkg/response"
)

// JobHandler 作业处理器
type JobHandler struct {
	service service.WorkflowService
}

// NewJobHandler 创建作业处理器
func NewJobHandler(service service.WorkflowService) *JobHandler {
	return &JobHandler{service: service}
}

// CreateJobRequest 创建作业请求
type CreateJobRequest struct {
	FlowID    int64  `json:"flow_id"`
	JobName   string `json:"job_name"`
	CreatedBy int64  `json:"created_by"`
}

// CreateJob 创建作业
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	job, err := h.service.CreateJob(req.FlowID, req.JobName, req.CreatedBy)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Created(w, job)
}

// GetJob 获取作业
func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid job id")
		return
	}

	job, jobTasks, err := h.service.GetJobWithTasks(id)
	if err != nil {
		response.NotFound(w, "job not found")
		return
	}

	result := map[string]interface{}{
		"job":       job,
		"job_tasks": jobTasks,
	}

	response.Success(w, result)
}

// ListJobs 获取作业列表
func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 20
	}

	jobs, err := h.service.ListJobs(limit, offset)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, jobs)
}

// StartJobRequest 启动作业请求
type StartJobRequest struct {
	JobID int64 `json:"job_id"`
}

// StartJob 启动作业
func (h *JobHandler) StartJob(w http.ResponseWriter, r *http.Request) {
	var req StartJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.service.StartJob(req.JobID); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, map[string]string{"message": "job started successfully"})
}

// StartTaskRequest 开始任务请求
type StartTaskRequest struct {
	JobTaskID  int64 `json:"job_task_id"`
	ExecutorID int64 `json:"executor_id"`
}

// StartTask 开始任务
func (h *JobHandler) StartTask(w http.ResponseWriter, r *http.Request) {
	var req StartTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.service.StartTask(req.JobTaskID, req.ExecutorID); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, map[string]string{"message": "task started successfully"})
}

// CompleteTaskRequest 完成任务请求
type CompleteTaskRequest struct {
	JobTaskID int64              `json:"job_task_id"`
	Result    models.TaskResult  `json:"result"`
}

// CompleteTask 完成任务
func (h *JobHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	var req CompleteTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.service.CompleteTask(req.JobTaskID, req.Result); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, map[string]string{"message": "task completed successfully"})
}

// FailTaskRequest 任务失败请求
type FailTaskRequest struct {
	JobTaskID    int64  `json:"job_task_id"`
	ErrorMessage string `json:"error_message"`
}

// FailTask 任务失败
func (h *JobHandler) FailTask(w http.ResponseWriter, r *http.Request) {
	var req FailTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.service.FailTask(req.JobTaskID, req.ErrorMessage); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, map[string]string{"message": "task marked as failed"})
}

// SkipTaskRequest 跳过任务请求
type SkipTaskRequest struct {
	JobTaskID  int64 `json:"job_task_id"`
	OperatorID int64 `json:"operator_id"`
}

// SkipTask 跳过任务
func (h *JobHandler) SkipTask(w http.ResponseWriter, r *http.Request) {
	var req SkipTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.service.SkipTask(req.JobTaskID, req.OperatorID); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, map[string]string{"message": "task skipped successfully"})
}

// RollbackTaskRequest 打回任务请求
type RollbackTaskRequest struct {
	JobTaskID      int64 `json:"job_task_id"`
	OperatorID     int64 `json:"operator_id"`
	TargetSequence int   `json:"target_sequence"`
}

// RollbackTask 打回任务
func (h *JobHandler) RollbackTask(w http.ResponseWriter, r *http.Request) {
	var req RollbackTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.service.RollbackTask(req.JobTaskID, req.OperatorID, req.TargetSequence); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, map[string]string{"message": "task rolled back successfully"})
}

// GetNextTask 获取下一个待执行的任务
func (h *JobHandler) GetNextTask(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.ParseInt(r.URL.Query().Get("job_id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid job id")
		return
	}

	task, err := h.service.GetNextTask(jobID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	if task == nil {
		response.Success(w, map[string]string{"message": "no more tasks"})
		return
	}

	response.Success(w, task)
}
