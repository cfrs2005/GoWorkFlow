package repository

import (
	"database/sql"
	"fmt"

	"github.com/cfrs2005/GoWorkFlow/internal/models"
)

// TaskRepository 任务仓储接口
type TaskRepository interface {
	Create(task *models.Task) error
	GetByID(id int64) (*models.Task, error)
	List(limit, offset int) ([]models.Task, error)
	Update(task *models.Task) error
	Delete(id int64) error
	GetByIDs(ids []int64) ([]models.Task, error)
}

type taskRepository struct {
	db *sql.DB
}

// NewTaskRepository 创建任务仓储
func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

// Create 创建任务
func (r *taskRepository) Create(task *models.Task) error {
	query := `
		INSERT INTO tasks (name, description, task_type, config, is_active)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, task.Name, task.Description, task.TaskType, task.Config, task.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	task.ID = id
	return nil
}

// GetByID 根据ID获取任务
func (r *taskRepository) GetByID(id int64) (*models.Task, error) {
	query := `
		SELECT id, name, description, task_type, config, is_active, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`
	task := &models.Task{}
	err := r.db.QueryRow(query, id).Scan(
		&task.ID, &task.Name, &task.Description, &task.TaskType,
		&task.Config, &task.IsActive, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// List 获取任务列表
func (r *taskRepository) List(limit, offset int) ([]models.Task, error) {
	query := `
		SELECT id, name, description, task_type, config, is_active, created_at, updated_at
		FROM tasks
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID, &task.Name, &task.Description, &task.TaskType,
			&task.Config, &task.IsActive, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Update 更新任务
func (r *taskRepository) Update(task *models.Task) error {
	query := `
		UPDATE tasks
		SET name = ?, description = ?, task_type = ?, config = ?, is_active = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, task.Name, task.Description, task.TaskType, task.Config, task.IsActive, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// Delete 删除任务
func (r *taskRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// GetByIDs 根据ID列表获取任务
func (r *taskRepository) GetByIDs(ids []int64) ([]models.Task, error) {
	if len(ids) == 0 {
		return []models.Task{}, nil
	}

	query := `
		SELECT id, name, description, task_type, config, is_active, created_at, updated_at
		FROM tasks
		WHERE id IN (?` + repeatPlaceholder(len(ids)-1) + `)
	`

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by ids: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID, &task.Name, &task.Description, &task.TaskType,
			&task.Config, &task.IsActive, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// repeatPlaceholder 生成重复的占位符
func repeatPlaceholder(count int) string {
	if count <= 0 {
		return ""
	}
	result := ""
	for i := 0; i < count; i++ {
		result += ",?"
	}
	return result
}
