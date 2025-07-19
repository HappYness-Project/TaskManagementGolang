package route

const prefix = "tasks_"

const UnexpectedErrMsg = "unexpected error occurred."

const (
	TaskGetUnauthorized          = prefix + "get_unauthorized"
	TaskGetNotFound              = prefix + "get_not_found"
	TaskGetRateLimitedExceeded   = prefix + "get_rate_limited_exceeded"
	TaskGetServerError           = prefix + "get_server_error"
	TaskGetTaskContainerNotFound = prefix + "get_taskcontainer_not_found"

	TaskCreateInvalidInput   = prefix + "create_invalid_input"
	TaskCreateServerError    = prefix + "create_server_error"
	TaskUpdateServerError    = prefix + "update_server_error"
	TaskStatusDoneError      = prefix + "status_done_error"
	TaskUpdateImportantError = prefix + "update_important_error"

	TaskDeleteInvalidID         = prefix + "delete_invalid_order_id"
	TaskDeleteNotFound          = prefix + "delete_not_found"
	TaskDeleteRateLimitExceeded = prefix + "delete_rate_limit_exceeded"
	TaskDeleteServerError       = prefix + "delete_server_error"
)
