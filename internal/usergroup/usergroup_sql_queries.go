package usergroup

const (
	sqlGetAllUsergroups = `SELECT * FROM public.usergroup`
	sqlGetById          = `SELECT * FROM public.usergroup WHERE id = $1`
)
