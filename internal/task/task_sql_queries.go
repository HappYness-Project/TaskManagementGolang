package task

const (
	sqlGetAllTasks              = `SELECT * FROM public.task`
	sqlGetTaskById              = `SELECT * FROM public.task WHERE id = $1`
	sqlGetAllTasksByContainerId = `SELECT t.id, t.name, t.description, t.type, t.created_at, t.updated_at, t.target_date, t.priority, t.category, t.is_completed, t.is_important
									FROM public.task t
									JOIN public.taskcontainer_task tct
									ON t.id = tct.task_id
									WHERE taskcontainer_id = $1`
	sqlCreateTask = `INSERT INTO public.task(id, name, description,type, created_at, updated_at, target_date, priority, category, is_completed, is_important)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	sqlDeleteTaskFromJoinTable = `DELETE FROM public.taskcontainer_task WHERE task_id=$1`
	sqlDeleteTask              = `DELETE FROM public.task WHERE id=$1`
	sqlUpdateTaskDoneField     = `UPDATE public.task SET is_completed=$1 WHERE id = $2;`
)
