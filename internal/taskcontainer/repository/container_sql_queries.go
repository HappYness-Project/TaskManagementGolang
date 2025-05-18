package repository

const (
	sqlGetAllContainers       = `SELECT id,name,description,is_active,usergroup_id FROM public.taskcontainer`
	sqlGetById                = `SELECT id,name,description,is_active,usergroup_id FROM public.taskcontainer WHERE id = $1`
	sqlGetContainersByGroupId = `SELECT id,name,description,is_active,usergroup_id FROM public.taskcontainer WHERE usergroup_id = $1`
	sqlCreateContainer        = `INSERT INTO public.taskcontainer(id, name, description, is_active, activity_level, type, usergroup_id)
								VALUES ($1,$2,$3,$4,$5,$6,$7);`
	sqlDeleteContainer              = `DELETE FROM public.taskcontainer WHERE id = $1;`
	sqlDeleteContainerByUsergroupId = `DELETE FROM taskcontainer WHERE usergroup_id = $1;`
)
