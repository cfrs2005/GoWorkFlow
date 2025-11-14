package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cfrs2005/GoWorkFlow/internal/config"
	"github.com/cfrs2005/GoWorkFlow/internal/engine"
	"github.com/cfrs2005/GoWorkFlow/internal/handler"
	"github.com/cfrs2005/GoWorkFlow/internal/repository"
	"github.com/cfrs2005/GoWorkFlow/internal/service"
	"github.com/cfrs2005/GoWorkFlow/pkg/database"
	"github.com/cfrs2005/GoWorkFlow/pkg/logger"
)

func main() {
	// 加载配置
	cfg := config.Load()
	logger.Infof("Starting workflow API server...")

	// 连接数据库
	db, err := database.NewDB(database.Config{
		DSN:             cfg.Database.GetDSN(),
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Info("Database connected successfully")

	// 初始化仓储层
	taskRepo := repository.NewTaskRepository(db.DB)
	flowRepo := repository.NewFlowRepository(db.DB)
	flowTaskRepo := repository.NewFlowTaskRepository(db.DB)
	jobRepo := repository.NewJobRepository(db.DB)
	jobTaskRepo := repository.NewJobTaskRepository(db.DB)

	// 初始化工作流引擎
	workflowEngine := engine.NewWorkflowEngine(
		db.DB,
		jobRepo,
		jobTaskRepo,
		flowRepo,
		flowTaskRepo,
	)

	// 初始化服务层
	workflowService := service.NewWorkflowService(
		db.DB,
		taskRepo,
		flowRepo,
		flowTaskRepo,
		jobRepo,
		jobTaskRepo,
		workflowEngine,
	)

	// 设置路由
	router := handler.NewRouter(workflowService)
	mux := router.Setup()

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Server listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
