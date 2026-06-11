package apperror

import "fmt"

type Code string

const (
	CodeBadRequest   Code = "bad_request"
	CodeNotFound     Code = "not_found"
	CodeExpired      Code = "expired"
	CodeRevoked      Code = "revoked"
	CodeForbidden    Code = "forbidden"
	CodeConflict     Code = "conflict"
	CodeInternal     Code = "internal"
	CodeUnauthorized Code = "unauthorized"
)

type Error struct {
	Code    Code
	Message string
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}
