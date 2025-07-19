package route

const prefix = "users_"

const (
	UserServerError = prefix + "server_error"
	UserDomainError = prefix + "domain_validation_error"

	UserGetNotFound            = prefix + "get_not_found"
	GetUserGroupsNotFound      = prefix + "get_usergroup_not_found"
	UserGetRateLimitedExceeded = prefix + "get_rate_limited_exceeded"

	UserCreateInvalidInput = prefix + "create_invalid_input"
	UserCreateServerError  = prefix + "create_server_error"

	UserUpdateServerError = prefix + "update_server_error"
)
