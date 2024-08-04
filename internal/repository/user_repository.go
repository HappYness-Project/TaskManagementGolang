package repository

import (
	"database/sql"

	"example.com/taskapp/internal/models"
)

type UserRepository interface {
	GetById(id int) (*models.User, error)
	GetUsersByGroupId(groupid int) ([]*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
}

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

// TODO
// func (r *UserRepo) GetById(id int) (*models.User, error) {
// 	return nil
// }
