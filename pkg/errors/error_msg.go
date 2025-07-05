package errors

const (
	// Token
	UnExpectedError     = "Expected error"
	ClaimsNotFound      = "Claims not found"
	TokenRequired       = "token required"
	TokenExpired        = "token expired"
	TokenInvalid        = "token invalid"
	InvalidRefreshToken = "invalid refresh token"

	PermissionDenied = "Permission denied"

	// DB
	RecordNotFound           = "record not found"
	BeginTransactionFailure  = "begin transaction failure"
	CommitTransactionFailure = "commit transaction failure"
	FailedDuringTransaction  = "failed during transaction"
	FailedToScanRow          = "failed to scan row"
	QueryExecutionFailure    = "query execution failure"
)
