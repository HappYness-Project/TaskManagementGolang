package taskcontainer

const (
	sqlGetAllContainers = `SELECT id,name,description,is_active,usergroup_id FROM public.taskcontainer`
	sqlGetById          = `SELECT id,name,description,is_active,usergroup_id FROM public.taskcontainer WHERE id = $1`
)
