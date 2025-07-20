package errors

const (
	ClaimsNotFound      = "Claims not found"
	TokenRequired       = "Token required"
	TokenExpired        = "Token expired"
	TokenInvalid        = "Token invalid"
	InvalidRefreshToken = "Invalid refresh token"

	InternalServerError = "Internal server error"
	Badrequest          = "Bad request"
	InvalidJsonBody     = "Invalid json body format"

	PermissionDenied = "Permission denied"

	// DB
	RecordNotFound           = "record not found"
	BeginTransactionFailure  = "begin transaction failure"
	CommitTransactionFailure = "commit transaction failure"
	FailedDuringTransaction  = "failed during transaction"
	FailedToScanRow          = "failed to scan row"
	QueryExecutionFailure    = "query execution failure"
)
