package user

const (
	sqlGetAllUsers       = `SELECT * FROM public.user`
	sqlGetUserById       = `SELECT * FROM public.user where id = $1`
	sqlGetUserByEmail    = `SELECT * FROM public.user where email = $1`
	sqlGetUsersByGroupId = `select * from public.user u
							INNER JOIN public.usergroup_user ugu
							ON u.id = ugu.user_id
							where ugu.usergroup_id = $1`
)
