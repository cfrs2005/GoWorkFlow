package response

import (
	"encoding/json"
	"net/http"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSON 返回 JSON 响应
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Success 成功响应
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Created 创建成功响应
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

// Error 错误响应
func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, Response{
		Code:    statusCode,
		Message: message,
	})
}

// BadRequest 400 错误请求
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

// NotFound 404 未找到
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

// InternalServerError 500 服务器错误
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}
