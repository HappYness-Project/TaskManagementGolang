package models

import "time"

type Task struct {
	TaskId      string    `json:"id"`
	TaskName    string    `json:task_name`
	TaskDesc    string    `json:task_desc`
	CreatedDate time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	TargetDate  time.Time `json:target_date`
	IsCompleted bool      `json:is_completed`
}
