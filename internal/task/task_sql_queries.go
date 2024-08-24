package task

const (
	sqlGetAllTasks              = `SELECT * FROM public.task`
	sqlGetTaskById              = `SELECT * FROM public.task WHERE id = $1`
	sqlGetAllTasksByContainerId = `SELECT t.id, t.name, t.description, t.type, t.created_at, t.updated_at, t.target_date, t.priority, t.category, t.is_completed, t.is_important
									FROM public.task t
									JOIN public.taskcontainer_task tct
									ON t.id = tct.task_id
									WHERE taskcontainer_id = $1`

	sqlDeleteTaskFromJoinTable = `DELETE FROM public.taskcontainer_task WHERE task_id=$1`
	sqlDeleteTask              = `DELETE FROM public.task WHERE id=$1`
	sqlUpdateTaskDoneField     = `UPDATE public.task SET is_completed=$1 WHERE id = $2;`
)
