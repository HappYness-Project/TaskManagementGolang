package task

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const dbTimeout = time.Second * 5

type TaskRepository interface {
	GetAllTasks() ([]*Task, error)
	GetTaskById(id string) (*Task, error)
	GetTasksByContainerId(containerId string) ([]*Task, error)
	CreateTask(task Task) (*Task, error)
	UpdateTask(task Task) error
	UpdateImportantTask(id string) error
	DeleteTask(id string) error
	DoneTask(id string) error
}
type TaskRepo struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepo {
	return &TaskRepo{
		DB: db,
	}
}

func (m *TaskRepo) GetAllTasks() ([]*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, sqlGetAllTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []*Task
	for rows.Next() {
		task, err := scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (m *TaskRepo) GetTaskById(id string) (*Task, error) {
	rows, err := m.DB.Query(sqlGetTaskById, id)
	if err != nil {
		return nil, err
	}

	task := new(Task)
	for rows.Next() {
		task, err = scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
	}
	return task, err
}

func (m *TaskRepo) GetTasksByContainerId(containerId string) ([]*Task, error) {
	rows, err := m.DB.Query(sqlGetAllTasksByContainerId, containerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task, err := scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (m *TaskRepo) CreateTask(task Task) (*Task, error) {
	_, err := m.DB.Exec(`INSERT INTO public.task(id, name, description, type, created_at, updated_at, target_date, priority, category, is_completed, is_important)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		task.TaskId, task.TaskName, task.TaskDesc, task.TaskType, task.CreatedAt, task.UpdatedAt, task.TargetDate, task.Priority, task.Category, task.IsCompleted, task.IsImportant)
	if err != nil {
		return nil, fmt.Errorf("unable to insert row : %w", err)
	}
	return &task, nil
}

func (m *TaskRepo) UpdateTask(task Task) error {
	return nil
}

func (m *TaskRepo) DeleteTask(id string) error {
	return nil
}

func (m *TaskRepo) DoneTask(id string) error {
	return nil
}

func (m *TaskRepo) UpdateImportantTask(id string) error {
	return nil
}

func scanRowsIntoTask(rows *sql.Rows) (*Task, error) {
	task := new(Task)
	err := rows.Scan(
		&task.TaskId,
		&task.TaskName,
		&task.TaskDesc,
		&task.TaskType,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.TargetDate,
		&task.Priority,
		&task.Category,
		&task.IsCompleted,
		&task.IsImportant,
	)
	if err != nil {
		return nil, err
	}

	return task, nil
}
