package usergroup

const (
	sqlGetAllUsergroups      = `SELECT * FROM public.usergroup`
	sqlGetById               = `SELECT * FROM public.usergroup WHERE id = $1`
	sqlGetUserGroupsByUserId = `SELECT ug.id, ug.name, ug.description, ug.type, ug.thumbnailurl, ug.is_active
								FROM public.usergroup ug
								INNER JOIN public.usergroup_user ugu
								ON ug.id = ugu.usergroup_id
								WHERE ugu.user_id = $1`

	sqlCreateUserGroup = `INSERT INTO public.usergroup(name, description, type, thumbnailurl, is_active)
							VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	sqlAddUserToUserGroup      = `INSERT INTO public.usergroup_user(usergroup_id, user_id) VALUES ($1, $2)`
	sqlRemoveUserFromUserGroup = `DELETE FROM public.usergroup_user WHERE usergroup_id = $1 AND user_id = $2`
)
