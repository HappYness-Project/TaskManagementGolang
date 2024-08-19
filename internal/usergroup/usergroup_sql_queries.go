package usergroup

const (
	sqlGetAllUsergroups      = `SELECT * FROM public.usergroup`
	sqlGetById               = `SELECT * FROM public.usergroup WHERE id = $1`
	sqlGetUserGroupsByUserId = `SELECT ug.id, ug.name, ug.description, ug.type, ug.thumbnailurl, ug.is_active
								FROM public.usergroup ug
								INNER JOIN public.usergroup_user ugu
								ON ug.id = ugu.user_id
								WHERE ugu.user_id = $1`
)
