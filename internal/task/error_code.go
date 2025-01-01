package task

const prefix = "tasks_"

const UnexpectedErrMsg = "unexpected error occurred."

const (
	TaskGetUnauthorized        = prefix + "get_unauthorized"
	TaskGetNotFound            = prefix + "get_not_found"
	TaskGetRateLimitedExceeded = prefix + "get_rate_limited_exceeded"
	TaskGetServerError         = prefix + "get_server_error"

	TaskCreateInvalidInput = prefix + "create_invalid_input"
	TaskCreateUnauthorized = prefix + "create_unauthorized"
	TaskCreateServerError  = prefix + "create_server_error"

	TaskDeleteInvalidID         = prefix + "delete_invalid_order_id"
	TaskDeleteUnauthorized      = prefix + "delete_unauthorized"
	TaskDeleteNotFound          = prefix + "delete_not_found"
	TaskDeleteRateLimitExceeded = prefix + "delete_rate_limit_exceeded"
	TaskDeleteServerError       = prefix + "delete_server_error"
)
