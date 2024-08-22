package task

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const dbTimeout = time.Second * 5

type TaskRepository interface {
	GetAllTasks() ([]*Task, error)
	GetTaskById(id string) (*Task, error)
	GetTasksByContainerId(containerId string) ([]*Task, error)
	CreateTask(taskcontainerId string, req CreateTaskDto) (int64, error)
	UpdateTask(req UpdateTaskDto) error
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
	defer rows.Close()

	var task *Task
	for rows.Next() {
		task, err = scanRowsIntoTask(rows)
		if err != nil {
			return nil, err
		}
	}
	return task, nil
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

func (m *TaskRepo) CreateTask(containerId string, req CreateTaskDto) (int64, error) {
	taskId := uuid.New()
	_, err := m.DB.Exec(`INSERT INTO public.task(id, name, description, created_at, updated_at, target_date, priority, category, is_completed, is_important)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		taskId, req.TaskName, req.TaskDesc, time.Now(), time.Now(), req.TargetDate, req.Priority, req.Category, false, false)
	if err != nil {
		return 0, fmt.Errorf("unable to insert into task table : %w", err)
	}
	_, err = m.DB.Exec(`INSERT INTO public.taskcontainer_task(taskcontainer_id, task_id)
		VALUES ($1, $2)`,
		containerId, taskId)

	if err != nil {
		return 0, fmt.Errorf("unable to insert into taskcontainer_task table : %w", err)
	}

	return 0, nil
}

func (m *TaskRepo) UpdateTask(task UpdateTaskDto) error {
	return nil
}

func (m *TaskRepo) DeleteTask(id string) error {
	_, err := m.DB.Exec(`DELETE FROM public.taskcontainer_task WHERE task_id=$1`, id)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec(`DELETE FROM public.task WHERE id=$1`, id)
	if err != nil {
		return err
	}
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
