package user

const (
	sqlGetAllUsers = `SELECT * FROM public.user`
	sqlGetUserById = `SELECT id, username, first_name, last_name, email, is_active,created_at,updated_at
						 FROM public.user
						 WHERE id = $1`
	sqlGetUserByEmail = `SELECT id, username, first_name, last_name, email, is_active,created_at,updated_at
							FROM public.user
							WHERE email = $1`
	sqlGetUserByUsername = `SELECT id, username, first_name, last_name, email, is_active,created_at,updated_at
							FROM public.user
							WHERE username = $1`
	sqlGetUsersByGroupId = `select * from public.user u
							INNER JOIN public.usergroup_user ugu
							ON u.id = ugu.user_id
							WHERE ugu.usergroup_id = $1`
)
