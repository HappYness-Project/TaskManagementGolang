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

	query := `select * from public.task`
	rows, err := m.DB.QueryContext(ctx, query)
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
	rows, err := m.DB.Query("SELECT * FROM public.task WHERE id = $1", id)
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT t.id, t.name, t.description, t.type, t.created_at, t.updated_at, t.target_date, t.priority, t.category, t.is_completed, t.is_important from public.task t
	 INNER JOIN public.taskcontainer_task tct
	 ON t.id = tct.task_id
	 WHERE tct.taskcontainer_id = $1`
	rows, err := m.DB.QueryContext(ctx, query, containerId)
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
	_, err := m.DB.Exec(`INSERT INTO public.task ()`)
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
