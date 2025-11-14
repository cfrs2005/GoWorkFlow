package repository

import (
	"database/sql"
	"fmt"

	"github.com/cfrs2005/GoWorkFlow/internal/models"
)

// FlowTaskRepository 流程任务仓储接口
type FlowTaskRepository interface {
	Create(flowTask *models.FlowTask) error
	GetByID(id int64) (*models.FlowTask, error)
	GetByFlowID(flowID int64) ([]models.FlowTask, error)
	Update(flowTask *models.FlowTask) error
	Delete(id int64) error
	DeleteByFlowID(flowID int64) error
}

type flowTaskRepository struct {
	db *sql.DB
}

// NewFlowTaskRepository 创建流程任务仓储
func NewFlowTaskRepository(db *sql.DB) FlowTaskRepository {
	return &flowTaskRepository{db: db}
}

// Create 创建流程任务
func (r *flowTaskRepository) Create(flowTask *models.FlowTask) error {
	query := `
		INSERT INTO flow_tasks (flow_id, task_id, sequence, is_optional, allow_rollback, condition_config)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		flowTask.FlowID, flowTask.TaskID, flowTask.Sequence,
		flowTask.IsOptional, flowTask.AllowRollback, flowTask.ConditionConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to create flow task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	flowTask.ID = id
	return nil
}

// GetByID 根据ID获取流程任务
func (r *flowTaskRepository) GetByID(id int64) (*models.FlowTask, error) {
	query := `
		SELECT id, flow_id, task_id, sequence, is_optional, allow_rollback, condition_config, created_at, updated_at
		FROM flow_tasks
		WHERE id = ?
	`
	flowTask := &models.FlowTask{}
	err := r.db.QueryRow(query, id).Scan(
		&flowTask.ID, &flowTask.FlowID, &flowTask.TaskID, &flowTask.Sequence,
		&flowTask.IsOptional, &flowTask.AllowRollback, &flowTask.ConditionConfig,
		&flowTask.CreatedAt, &flowTask.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("flow task not found")
		}
		return nil, fmt.Errorf("failed to get flow task: %w", err)
	}

	return flowTask, nil
}

// GetByFlowID 根据流程ID获取所有流程任务
func (r *flowTaskRepository) GetByFlowID(flowID int64) ([]models.FlowTask, error) {
	query := `
		SELECT id, flow_id, task_id, sequence, is_optional, allow_rollback, condition_config, created_at, updated_at
		FROM flow_tasks
		WHERE flow_id = ?
		ORDER BY sequence ASC
	`
	rows, err := r.db.Query(query, flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get flow tasks: %w", err)
	}
	defer rows.Close()

	var flowTasks []models.FlowTask
	for rows.Next() {
		var flowTask models.FlowTask
		if err := rows.Scan(
			&flowTask.ID, &flowTask.FlowID, &flowTask.TaskID, &flowTask.Sequence,
			&flowTask.IsOptional, &flowTask.AllowRollback, &flowTask.ConditionConfig,
			&flowTask.CreatedAt, &flowTask.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan flow task: %w", err)
		}
		flowTasks = append(flowTasks, flowTask)
	}

	return flowTasks, nil
}

// Update 更新流程任务
func (r *flowTaskRepository) Update(flowTask *models.FlowTask) error {
	query := `
		UPDATE flow_tasks
		SET task_id = ?, sequence = ?, is_optional = ?, allow_rollback = ?, condition_config = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		flowTask.TaskID, flowTask.Sequence, flowTask.IsOptional,
		flowTask.AllowRollback, flowTask.ConditionConfig, flowTask.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update flow task: %w", err)
	}

	return nil
}

// Delete 删除流程任务
func (r *flowTaskRepository) Delete(id int64) error {
	query := `DELETE FROM flow_tasks WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete flow task: %w", err)
	}

	return nil
}

// DeleteByFlowID 删除流程的所有任务
func (r *flowTaskRepository) DeleteByFlowID(flowID int64) error {
	query := `DELETE FROM flow_tasks WHERE flow_id = ?`
	_, err := r.db.Exec(query, flowID)
	if err != nil {
		return fmt.Errorf("failed to delete flow tasks: %w", err)
	}

	return nil
}
