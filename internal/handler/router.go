package handler

import (
	"net/http"

	"github.com/cfrs2005/GoWorkFlow/internal/service"
)

// Router 路由器
type Router struct {
	taskHandler *TaskHandler
	flowHandler *FlowHandler
	jobHandler  *JobHandler
}

// NewRouter 创建路由器
func NewRouter(service service.WorkflowService) *Router {
	return &Router{
		taskHandler: NewTaskHandler(service),
		flowHandler: NewFlowHandler(service),
		jobHandler:  NewJobHandler(service),
	}
}

// Setup 设置路由
func (router *Router) Setup() *http.ServeMux {
	mux := http.NewServeMux()

	// Task 路由
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			router.taskHandler.CreateTask(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				router.taskHandler.GetTask(w, r)
			} else {
				router.taskHandler.ListTasks(w, r)
			}
		case http.MethodPut:
			router.taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			router.taskHandler.DeleteTask(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Flow 路由
	mux.HandleFunc("/api/flows", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			router.flowHandler.CreateFlow(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				router.flowHandler.GetFlow(w, r)
			} else {
				router.flowHandler.ListFlows(w, r)
			}
		case http.MethodPut:
			router.flowHandler.UpdateFlow(w, r)
		case http.MethodDelete:
			router.flowHandler.DeleteFlow(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Job 路由
	mux.HandleFunc("/api/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			router.jobHandler.CreateJob(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				router.jobHandler.GetJob(w, r)
			} else {
				router.jobHandler.ListJobs(w, r)
			}
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/jobs/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		router.jobHandler.StartJob(w, r)
	})

	mux.HandleFunc("/api/jobs/next-task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		router.jobHandler.GetNextTask(w, r)
	})

	// JobTask 路由
	mux.HandleFunc("/api/tasks/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		router.jobHandler.StartTask(w, r)
	})

	mux.HandleFunc("/api/tasks/complete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		router.jobHandler.CompleteTask(w, r)
	})

	mux.HandleFunc("/api/tasks/fail", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		router.jobHandler.FailTask(w, r)
	})

	mux.HandleFunc("/api/tasks/skip", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		router.jobHandler.SkipTask(w, r)
	})

	mux.HandleFunc("/api/tasks/rollback", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		router.jobHandler.RollbackTask(w, r)
	})

	// 健康检查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return mux
}
