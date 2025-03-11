package model

import (
	"time"
)

type User struct {
	Id             int       `json:"id"`
	UserId         string    `json:"user_id"`
	UserName       string    `json:"username"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DefaultGroupId int       `json:"default_group_id"`
}
