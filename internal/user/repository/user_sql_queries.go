package repository

const (
	sqlGetAllUsers     = `SELECT * FROM public.user`
	sqlGetUserByUserId = `SELECT id, user_id, username, first_name, last_name, email, is_active,created_at,updated_at, default_group_id
						 FROM public.user
						 WHERE user_id = $1`
	sqlGetUserByEmail = `SELECT id, user_id, username, first_name, last_name, email, is_active,created_at,updated_at, default_group_id
							FROM public.user
							WHERE email = $1`
	sqlGetUserByUsername = `SELECT id, user_id, username, first_name, last_name, email, is_active,created_at,updated_at, default_group_id
							FROM public.user
							WHERE username = $1`
	sqlGetUsersByGroupId = `SELECT id, u.user_id, username, first_name, last_name, email, is_active,created_at,updated_at, default_group_id from public.user u
							INNER JOIN public.usergroup_user ugu
							ON u.id = ugu.user_id
							WHERE ugu.usergroup_id = $1`
	sqlCreateUser        = `INSERT INTO public.user(user_id, username, first_name, last_name, email, is_active, created_at, updated_at, default_group_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	sqlCreateUserSetting = `INSERT INTO public.usersetting VALUES($1, $2)`
	sqlUpdateUser        = `UPDATE public.user SET first_name=$2, last_name=$3, email=$4, default_group_id=$5, updated_at=$6 WHERE id=$1`
)
