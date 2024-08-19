package user

import (
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/user"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup"
)

type UserDetailDto struct {
	Id          int                   `json:"id"`
	UserName    string                `json:"username"`
	FirstName   string                `json:"first_name"`
	LastName    string                `json:"last_name"`
	Email       string                `json:"email"`
	IsActive    bool                  `json:"is_active"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	UserSetting user.UserSetting      `json:"usersetting"`
	UserGroup   []usergroup.UserGroup `json:"usergroup"`
}
