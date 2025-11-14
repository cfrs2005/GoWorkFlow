# Router Update Instructions

需要手动更新 `internal/handler/router.go` 文件，添加以下内容：

## 1. 更新 Router 结构体

```go
// Router 路由器
type Router struct {
	taskHandler         *TaskHandler
	flowHandler         *FlowHandler
	jobHandler          *JobHandler
	jobContextHandler   *JobContextHandler
	executorHandler     *ExecutorHandler  // 新增
}
```

## 2. 更新 NewRouter 函数

```go
// NewRouter 创建路由器
func NewRouter(
	service service.WorkflowService,
	jobContextRepo repository.JobContextRepository,
	taskExecutorService *service.TaskExecutorService,  // 新增参数
) *Router {
	return &Router{
		taskHandler:       NewTaskHandler(service),
		flowHandler:       NewFlowHandler(service),
		jobHandler:        NewJobHandler(service),
		jobContextHandler: NewJobContextHandler(jobContextRepo),
		executorHandler:   NewExecutorHandler(taskExecutorService),  // 新增
	}
}
```

## 3. 在 Setup() 方法的 mux.HandleFunc("/api/jobs/start", ...) 之后添加新路由

```go
	// 自动执行路由（添加在 /api/jobs/start 之后）
	mux.HandleFunc("/api/jobs/auto-execute", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		router.executorHandler.AutoExecuteJob(w, r)
	})

	// 任务执行路由（添加在 /api/tasks/execute 位置）
	mux.HandleFunc("/api/tasks/execute", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		router.executorHandler.ExecuteTask(w, r)
	})
```

## 或者，直接替换整个 router.go 文件

如果手动编辑麻烦，可以运行以下命令（已提供完整的新文件在下方）：

```bash
# 备份原文件
cp internal/handler/router.go internal/handler/router.go.backup

# 使用新文件（见下方完整代码）
```
