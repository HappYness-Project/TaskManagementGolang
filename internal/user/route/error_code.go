package route

import "errors"

const prefix = "users_"

const UnexpectedErrMsg = "unexpected error occurred."

const (
	UserGetUnauthorized        = prefix + "get_unauthorized"
	UserGetNotFound            = prefix + "get_not_found"
	UserGetRateLimitedExceeded = prefix + "get_rate_limited_exceeded"
	UserGetServerError         = prefix + "get_server_error"

	UserCreateInvalidInput = prefix + "create_invalid_input"
	UserCreateUnauthorized = prefix + "create_unauthorized"
	UserCreateServerError  = prefix + "create_server_error"
)

var (
	CreateUserError = errors.New("cannot create a user")
)
