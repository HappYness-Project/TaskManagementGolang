package repository

import (
	"database/sql"

	"example.com/taskapp/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllTaskContainers() ([]*models.TaskContainer, error)
}
