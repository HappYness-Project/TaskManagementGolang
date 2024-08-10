package user

import "time"

type User struct {
	Id            int       `json:"id"`
	UserName      string    `json:user_name`
	FirstName     string    `json:first_name`
	LastName      string    `json:last_name`
	Email         string    `json:email`
	IsActive      bool      `json:is_active`
	CreatedAt     time.Time `json:created_at`
	UserSettingId int       `json:usersetting_id`
}
