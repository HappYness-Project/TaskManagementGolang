package repository

import "example.com/taskapp/internal/models"

type UserRepository interface {
	GetById(id int) (*models.User, error)
	GetUsersByGroupId(groupid int) ([]*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
}
