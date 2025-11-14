package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cfrs2005/GoWorkFlow/internal/models"
	"github.com/cfrs2005/GoWorkFlow/internal/service"
	"github.com/cfrs2005/GoWorkFlow/pkg/response"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	service service.WorkflowService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(service service.WorkflowService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask 创建任务
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	task.IsActive = true
	if err := h.service.CreateTask(&task); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Created(w, task)
}

// GetTask 获取任务
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid task id")
		return
	}

	task, err := h.service.GetTask(id)
	if err != nil {
		response.NotFound(w, "task not found")
		return
	}

	response.Success(w, task)
}

// ListTasks 获取任务列表
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 20
	}

	tasks, err := h.service.ListTasks(limit, offset)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, tasks)
}

// UpdateTask 更新任务
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.service.UpdateTask(&task); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, task)
}

// DeleteTask 删除任务
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid task id")
		return
	}

	if err := h.service.DeleteTask(id); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, map[string]string{"message": "task deleted successfully"})
}
