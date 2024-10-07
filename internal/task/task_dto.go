package task

import "time"

type CreateTaskDto struct {
	TaskName   string    `json:"name"`
	TaskDesc   string    `json:"description"`
	TargetDate time.Time `json:"target_date"`
	Priority   string    `json:"priority"`
	Category   string    `json:"category"`
}

type UpdateTaskDto struct {
	TaskName   string    `json:"name"`
	TaskDesc   string    `json:"description"`
	TargetDate time.Time `json:"target_date"`
	Priority   string    `json:"priority"`
	Category   string    `json:"category"`
}
