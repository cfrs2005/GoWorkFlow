package repository

import (
	"database/sql"
	"fmt"

	"github.com/cfrs2005/GoWorkFlow/internal/models"
)

// FlowRepository 流程仓储接口
type FlowRepository interface {
	Create(flow *models.Flow) error
	GetByID(id int64) (*models.Flow, error)
	List(limit, offset int) ([]models.Flow, error)
	Update(flow *models.Flow) error
	Delete(id int64) error
	GetFlowWithTasks(flowID int64) (*models.Flow, []models.FlowTask, error)
}

type flowRepository struct {
	db *sql.DB
}

// NewFlowRepository 创建流程仓储
func NewFlowRepository(db *sql.DB) FlowRepository {
	return &flowRepository{db: db}
}

// Create 创建流程
func (r *flowRepository) Create(flow *models.Flow) error {
	query := `
		INSERT INTO flows (name, description, version, is_active, created_by)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, flow.Name, flow.Description, flow.Version, flow.IsActive, flow.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to create flow: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	flow.ID = id
	return nil
}

// GetByID 根据ID获取流程
func (r *flowRepository) GetByID(id int64) (*models.Flow, error) {
	query := `
		SELECT id, name, description, version, is_active, created_by, created_at, updated_at
		FROM flows
		WHERE id = ?
	`
	flow := &models.Flow{}
	err := r.db.QueryRow(query, id).Scan(
		&flow.ID, &flow.Name, &flow.Description, &flow.Version,
		&flow.IsActive, &flow.CreatedBy, &flow.CreatedAt, &flow.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("flow not found")
		}
		return nil, fmt.Errorf("failed to get flow: %w", err)
	}

	return flow, nil
}

// List 获取流程列表
func (r *flowRepository) List(limit, offset int) ([]models.Flow, error) {
	query := `
		SELECT id, name, description, version, is_active, created_by, created_at, updated_at
		FROM flows
		WHERE is_active = 1
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list flows: %w", err)
	}
	defer rows.Close()

	var flows []models.Flow
	for rows.Next() {
		var flow models.Flow
		if err := rows.Scan(
			&flow.ID, &flow.Name, &flow.Description, &flow.Version,
			&flow.IsActive, &flow.CreatedBy, &flow.CreatedAt, &flow.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan flow: %w", err)
		}
		flows = append(flows, flow)
	}

	return flows, nil
}

// Update 更新流程
func (r *flowRepository) Update(flow *models.Flow) error {
	query := `
		UPDATE flows
		SET name = ?, description = ?, version = ?, is_active = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, flow.Name, flow.Description, flow.Version, flow.IsActive, flow.ID)
	if err != nil {
		return fmt.Errorf("failed to update flow: %w", err)
	}

	return nil
}

// Delete 删除流程
func (r *flowRepository) Delete(id int64) error {
	query := `DELETE FROM flows WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete flow: %w", err)
	}

	return nil
}

// GetFlowWithTasks 获取流程及其关联的任务
func (r *flowRepository) GetFlowWithTasks(flowID int64) (*models.Flow, []models.FlowTask, error) {
	// 获取流程
	flow, err := r.GetByID(flowID)
	if err != nil {
		return nil, nil, err
	}

	// 获取流程任务
	query := `
		SELECT ft.id, ft.flow_id, ft.task_id, ft.sequence, ft.is_optional,
		       ft.allow_rollback, ft.condition_config, ft.created_at, ft.updated_at,
		       t.id, t.name, t.description, t.task_type, t.config, t.is_active,
		       t.created_at, t.updated_at
		FROM flow_tasks ft
		INNER JOIN tasks t ON ft.task_id = t.id
		WHERE ft.flow_id = ?
		ORDER BY ft.sequence ASC
	`
	rows, err := r.db.Query(query, flowID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get flow tasks: %w", err)
	}
	defer rows.Close()

	var flowTasks []models.FlowTask
	for rows.Next() {
		var flowTask models.FlowTask
		var task models.Task
		if err := rows.Scan(
			&flowTask.ID, &flowTask.FlowID, &flowTask.TaskID, &flowTask.Sequence,
			&flowTask.IsOptional, &flowTask.AllowRollback, &flowTask.ConditionConfig,
			&flowTask.CreatedAt, &flowTask.UpdatedAt,
			&task.ID, &task.Name, &task.Description, &task.TaskType,
			&task.Config, &task.IsActive, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, nil, fmt.Errorf("failed to scan flow task: %w", err)
		}
		flowTask.Task = &task
		flowTasks = append(flowTasks, flowTask)
	}

	return flow, flowTasks, nil
}
