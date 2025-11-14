package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cfrs2005/GoWorkFlow/internal/repository"
	"github.com/cfrs2005/GoWorkFlow/pkg/response"
)

// JobContextHandler 作业上下文处理器
type JobContextHandler struct {
	repo repository.JobContextRepository
}

// NewJobContextHandler 创建作业上下文处理器
func NewJobContextHandler(repo repository.JobContextRepository) *JobContextHandler {
	return &JobContextHandler{repo: repo}
}

// GetJobContext 获取作业上下文
// GET /api/jobs/{id}/context
func (h *JobContextHandler) GetJobContext(w http.ResponseWriter, r *http.Request) {
	// 从 URL 路径中提取 job_id
	jobID, err := h.extractJobID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	context, err := h.repo.GetByJobID(jobID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, context)
}

// UpdateJobContext 更新作业上下文
// PUT /api/jobs/{id}/context
func (h *JobContextHandler) UpdateJobContext(w http.ResponseWriter, r *http.Request) {
	// 从 URL 路径中提取 job_id
	jobID, err := h.extractJobID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	// 解析请求体
	var contextData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&contextData); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 更新每个键值对
	for key, value := range contextData {
		if err := h.repo.Set(jobID, key, value); err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.Success(w, map[string]string{"message": "Context updated successfully"})
}

// extractJobID 从请求路径中提取 job_id
func (h *JobContextHandler) extractJobID(r *http.Request) (int64, error) {
	// 路径格式: /api/jobs/{id}/context
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		return 0, strconv.ErrSyntax
	}

	return strconv.ParseInt(parts[2], 10, 64)
}
