package domain

import "fmt"

type ErrorCode string

const (
	ErrNotFound         ErrorCode = "NOT_FOUND"
	ErrInvalidQuery     ErrorCode = "INVALID_QUERY"
	ErrIndexingFailed   ErrorCode = "INDEXING_FAILED"
	ErrSearchFailed     ErrorCode = "SEARCH_FAILED"
	ErrInternalError    ErrorCode = "INTERNAL_ERROR"
	ErrConnectionFailed ErrorCode = "CONNECTION_FAILED"
)

type DomainError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewDomainError(code ErrorCode, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

func NewDomainErrorWithCause(code ErrorCode, message string, cause error) *DomainError {
	return &DomainError{Code: code, Message: message, Cause: cause}
}
