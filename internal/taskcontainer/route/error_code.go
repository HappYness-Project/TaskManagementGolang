package route

const prefix = "taskcontainers_"

const (
	TaskContainerServerError = prefix + "server_error"
	TaskContainerDomainError = prefix + "domain_validation_error"
	TaskContainerGetError    = prefix + "get_server_error"

	TaskContainerGetNotFound = prefix + "get_not_found"
	DeleteTaskContainerError = prefix + "delete_server_error"
	// UserGroupGetRateLimitedExceeded = prefix + "get_rate_limited_exceeded"
	// UserNotFound                    = prefix + "user_get_not_found"
	// UserGroupCreationFailure        = prefix + "create_error"
	// UserGroupAddUserError           = prefix + "add_user_error"
	// DeleteUserGroupError            = prefix + "delete_server_error"

	RemoveUserFromUserGroupError = prefix + "remove_user_from_Usergroup_server_error"
)
