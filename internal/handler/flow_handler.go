package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cfrs2005/GoWorkFlow/internal/models"
	"github.com/cfrs2005/GoWorkFlow/internal/service"
	"github.com/cfrs2005/GoWorkFlow/pkg/response"
)

// FlowHandler 流程处理器
type FlowHandler struct {
	service service.WorkflowService
}

// NewFlowHandler 创建流程处理器
func NewFlowHandler(service service.WorkflowService) *FlowHandler {
	return &FlowHandler{service: service}
}

// CreateFlowRequest 创建流程请求
type CreateFlowRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Version     string  `json:"version"`
	TaskIDs     []int64 `json:"task_ids"`
	CreatedBy   int64   `json:"created_by"`
}

// CreateFlow 创建流程
func (h *FlowHandler) CreateFlow(w http.ResponseWriter, r *http.Request) {
	var req CreateFlowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	flow := &models.Flow{
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		IsActive:    true,
		CreatedBy:   req.CreatedBy,
	}

	if err := h.service.CreateFlow(flow, req.TaskIDs); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Created(w, flow)
}

// GetFlow 获取流程
func (h *FlowHandler) GetFlow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid flow id")
		return
	}

	flow, flowTasks, err := h.service.GetFlowWithTasks(id)
	if err != nil {
		response.NotFound(w, "flow not found")
		return
	}

	result := map[string]interface{}{
		"flow":       flow,
		"flow_tasks": flowTasks,
	}

	response.Success(w, result)
}

// ListFlows 获取流程列表
func (h *FlowHandler) ListFlows(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 20
	}

	flows, err := h.service.ListFlows(limit, offset)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, flows)
}

// UpdateFlow 更新流程
func (h *FlowHandler) UpdateFlow(w http.ResponseWriter, r *http.Request) {
	var flow models.Flow
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.service.UpdateFlow(&flow); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, flow)
}

// DeleteFlow 删除流程
func (h *FlowHandler) DeleteFlow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		response.BadRequest(w, "invalid flow id")
		return
	}

	if err := h.service.DeleteFlow(id); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, map[string]string{"message": "flow deleted successfully"})
}
