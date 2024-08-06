package task

import "time"

type Task struct {
	TaskId      string    `json:"id"`
	TaskName    string    `json:"name"`
	TaskDesc    string    `json:"description"`
	TaskType    string    `json:"type"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	TargetDate  time.Time `json:"target_date"`
	Priority    string    `json:"priority"`
	Category    string    `json:"category"`
	IsCompleted bool      `json:"is_completed"`
	IsImportant bool      `json:"is_important"`
}
