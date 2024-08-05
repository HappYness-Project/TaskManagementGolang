package task

import "time"

type Task struct {
	TaskId      string    `json:"id"`
	TaskName    string    `json:task_name`
	TaskDesc    string    `json:task_desc`
	TaskType    string    `json:task_type`
	Priority    string    `json:task_priority`
	CreatedDate time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	TargetDate  time.Time `json:"-"`
	IsCompleted bool      `json:is_completed`
	IsImportant bool      `json:is_important`
}
