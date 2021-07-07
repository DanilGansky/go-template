package errors

type ErrorCode string

const (
	ValidationError ErrorCode = "validation_error"
	InternalError   ErrorCode = "internal_error"
	NotFoundError   ErrorCode = "not_found"
	DuplicateError  ErrorCode = "duplicate_error"
	Unauthorized    ErrorCode = "unauthorized"
)

func (code ErrorCode) String() string {
	return string(code)
}
