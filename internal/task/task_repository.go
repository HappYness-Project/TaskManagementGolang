package task

import (
	"context"
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 5

type TaskRepository interface {
	GetAllTasks() ([]*Task, error)
	GetTaskById(id string) (*Task, error)
	CreateTask(task Task) (*Task, error)
	UpdateTask(task Task) (Task, error)
	UpdateImportantTask(id string)
	DeleteTask(id string)
	DoneTask(id string)
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
		task, err := scanRowsIntoContainer(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func scanRowsIntoContainer(rows *sql.Rows) (*Task, error) {
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
