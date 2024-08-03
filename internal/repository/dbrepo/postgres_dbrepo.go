package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"example.com/taskapp/internal/models"
)

type PostgresDbRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 5

func (m *PostgresDbRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDbRepo) AllTaskContainers() ([]*models.TaskContainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select id,name,description from public.taskcontainer
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containers []*models.TaskContainer
	for rows.Next() {
		var container models.TaskContainer
		err := rows.Scan(
			&container.ContainerId,
			&container.ContainerName,
			&container.ContainerDesc,
		)
		if err != nil {
			return nil, err
		}

		containers = append(containers, &container)
	}
	return containers, nil

}
