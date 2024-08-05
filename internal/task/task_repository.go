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

	container := new(Task)
	for rows.Next() {
		container, err = scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
	}
	return container, err
}

func (m *TaskRepo) GetTasksByContainerId(containerId string) ([]*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select * from public.task` // TODO Update this statement.
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
	)
	if err != nil {
		return nil, err
	}

	return task, nil
}
