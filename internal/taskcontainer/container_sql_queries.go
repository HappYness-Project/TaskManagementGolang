package taskcontainer

const (
	sqlGetAllContainers = `SELECT id,name,description FROM public.taskcontainer`
	sqlGetById          = `SELECT id,name,description FROM public.taskcontainer WHERE id = $1`
)
