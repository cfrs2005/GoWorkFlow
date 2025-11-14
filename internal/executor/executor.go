package executor

import (
	"context"
	"fmt"
)

// Executor 任务执行器接口
type Executor interface {
	// Execute 执行任务
	// input: 任务输入参数
	// context: 作业上下文（用于任务间数据共享）
	// 返回: 输出结果和错误
	Execute(ctx context.Context, input map[string]interface{}, jobContext map[string]string) (map[string]interface{}, error)

	// Name 返回执行器名称
	Name() string
}

// ExecutorRegistry 执行器注册表
type ExecutorRegistry struct {
	executors map[string]Executor
}

// NewExecutorRegistry 创建执行器注册表
func NewExecutorRegistry() *ExecutorRegistry {
	return &ExecutorRegistry{
		executors: make(map[string]Executor),
	}
}

// Register 注册执行器
func (r *ExecutorRegistry) Register(executor Executor) {
	r.executors[executor.Name()] = executor
}

// Get 获取执行器
func (r *ExecutorRegistry) Get(name string) (Executor, error) {
	executor, ok := r.executors[name]
	if !ok {
		return nil, fmt.Errorf("executor not found: %s", name)
	}
	return executor, nil
}

// List 列出所有执行器
func (r *ExecutorRegistry) List() []string {
	names := make([]string, 0, len(r.executors))
	for name := range r.executors {
		names = append(names, name)
	}
	return names
}

// Global executor registry
var globalRegistry = NewExecutorRegistry()

// RegisterExecutor 注册全局执行器
func RegisterExecutor(executor Executor) {
	globalRegistry.Register(executor)
}

// GetExecutor 获取全局执行器
func GetExecutor(name string) (Executor, error) {
	return globalRegistry.Get(name)
}

// ListExecutors 列出所有执行器
func ListExecutors() []string {
	return globalRegistry.List()
}
