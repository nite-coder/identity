package domain

import "google.golang.org/grpc/codes"

// AppError handles application exception.
type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Status  codes.Code             `json:"status"`
	Details map[string]interface{} `json:"details"`
}

func (e AppError) Error() string {
	return e.Message
}

// New functions create a new AppError instance
func New(code, message string) AppError {
	return AppError{Code: code, Message: message}
}

var (
	ErrNotFound                    = &AppError{Code: "NOT_FOUND", Message: "resource was not found", Status: codes.NotFound}
	ErrWrongStatus                 = &AppError{Code: "WRONG_STATUS", Message: "status of resource is wrong", Status: codes.Internal}
	ErrStale                       = &AppError{Code: "STALE", Message: "resource is stale.  please retry", Status: codes.Internal}
	ErrInvalidInput                = &AppError{Code: "INVAID_INPUT", Message: "input is invalid.", Status: codes.InvalidArgument}
	ErrAlreadyExists               = &AppError{Code: "ALREADY_EXISTS", Message: "resource already exists", Status: codes.AlreadyExists}
	ErrUsernameOrPasswordIncorrect = &AppError{Code: "USERNAME_OR_PASSWORD_INCORRECT", Message: "username or paassword is incorrect", Status: codes.InvalidArgument}
	ErrAccountDisable              = &AppError{Code: "ACCOUNT_DISABLE", Message: "the account is disabled", Status: codes.Internal}
)
