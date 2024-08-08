package task

const (
	sqlGetAllTasks              = `SELECT * FROM public.task`
	sqlGetTaskById              = `SELECT * FROM public.task WHERE id = $1`
	sqlGetAllTasksByContainerId = `SELECT t.id, t.name, t.description, t.type, t.created_at, t.updated_at, t.target_date, t.priority, t.category, t.is_completed, t.is_important
									FROM public.task t
									JOIN public.taskcontainer_task tct
									ON t.id = tct.task_id
									WHERE taskcontainer_id = $1`
)
