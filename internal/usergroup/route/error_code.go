package route

const prefix = "usergroups_"

const (
	UserGroupServerError = prefix + "server_error"
	UserGroupDomainError = prefix + "domain_validation_error"

	UserGroupGetNotFound            = prefix + "get_not_found"
	UserGroupGetRateLimitedExceeded = prefix + "get_rate_limited_exceeded"

	UserNotFound = prefix + "user_get_not_found"

	UserGroupCreationFailure = prefix + "create_error"

	UserGroupAddUserError = prefix + "add_user_error"
	// UserCreateInvalidInput = prefix + "create_invalid_input"
	// UserCreateUnauthorized = prefix + "create_unauthorized"
	// UserCreateServerError  = prefix + "create_server_error"
	DeleteUserGroupError = prefix + "delete_server_error"

	RemoveUserFromUserGroupError = prefix + "remove_user_from_Usergroup_server_error"
)

// var (
// 	CreateUserError = errors.New("cannot create a user")
// )
