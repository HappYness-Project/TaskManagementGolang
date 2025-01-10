package route

import (
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/user/model"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup"
)

type CreateUserDto struct {
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
type UpdateUserDto struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
type UserDetailDto struct {
	Id          int                    `json:"id"`
	UserName    string                 `json:"username"`
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Email       string                 `json:"email"`
	IsActive    bool                   `json:"is_active"`
	UserSetting *model.UserSetting     `json:"user_setting"`
	UserGroup   []*usergroup.UserGroup `json:"user_groups"`
}
