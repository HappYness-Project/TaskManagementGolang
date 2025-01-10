package repository

const (
	sqlGetAllUsers = `SELECT * FROM public.user`
	sqlGetUserById = `SELECT id, username, first_name, last_name, email, is_active,created_at,updated_at, usersetting_id
						 FROM public.user
						 WHERE id = $1`
	sqlGetUserByEmail = `SELECT id, username, first_name, last_name, email, is_active,created_at,updated_at, usersetting_id
							FROM public.user
							WHERE email = $1`
	sqlGetUserByUsername = `SELECT id, username, first_name, last_name, email, is_active,created_at,updated_at, usersetting_id
							FROM public.user
							WHERE username = $1`
	sqlGetUsersByGroupId = `SELECT id, username, first_name, last_name, email, is_active,created_at,updated_at, usersetting_id from public.user u
							INNER JOIN public.usergroup_user ugu
							ON u.id = ugu.user_id
							WHERE ugu.usergroup_id = $1`
	sqlGetDefaultGroupId  = `SELECT default_group_id FROM public.usersetting WHERE id = $1`
	sqlGetUserSettingById = `SELECT * FROM public.usersetting WHERE id = $1`
	sqlCreateUser         = `INSERT INTO public.user(id, username, first_name, last_name, email, is_active, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	sqlCreateUserSetting  = `INSERT INTO public.usersetting VALUES($1, $2)`
	sqlUpdateUser         = `UPDATE public.user SET first_name=$2, last_name=$3, email=$4, updated_at=$5 WHERE id=$1`
)
