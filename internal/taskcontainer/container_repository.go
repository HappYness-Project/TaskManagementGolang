package taskcontainer

import (
	"context"
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 5

type ContainerRepository interface {
	AllTaskContainers() ([]*TaskContainer, error)
	GetById(id string) (*TaskContainer, error)
}

type ContainerRepo struct {
	DB *sql.DB
}

func NewContainerRepository(db *sql.DB) *ContainerRepo {
	return &ContainerRepo{
		DB: db,
	}
}

func (m *ContainerRepo) Connection() *sql.DB {
	return m.DB
}

func (m *ContainerRepo) AllTaskContainers() ([]*TaskContainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id,name,description from public.taskcontainer`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containers []*TaskContainer
	for rows.Next() {
		container, err := scanRowsIntoContainer(rows)
		if err != nil {
			return nil, err
		}

		containers = append(containers, &container)
	}
	return containers, nil
}

func (m *ContainerRepo) GetById(id string) (*TaskContainer, error) {
	rows, err := m.DB.Query("SELECT id,name,description FROM public.taskcontainer WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	container := new(TaskContainer)
	for rows.Next() {
		container, err = scanRowsIntoContainer(rows)
		if err != nil {
			return nil, err
		}
	}
	return container, err
}
func scanRowsIntoContainer(rows *sql.Rows) (*TaskContainer, error) {
	container := new(TaskContainer)
	err := rows.Scan(
		&container.ContainerId,
		&container.ContainerName,
		&container.ContainerDesc,
	)
	if err != nil {
		return nil, err
	}

	return container, nil
}
