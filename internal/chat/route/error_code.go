package route

const prefix = "chats_"
const (
	ChatGetNotFound        = prefix + "get_not_found"
	ChatGetServerError     = prefix + "get_server_error"
	ChatCreateInvalidInput = prefix + "create_invalid_input"
	ChatDeleteServerError  = prefix + "delete_server_error"
)
