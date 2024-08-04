package global

import "fmt"

const (
    CodeSuccess        = "200"
    CodeNoData         = "204"
    CodeInvalidRequest = "400"
    CodeAuthFailed     = "401"
    CodeExceededLimit  = "402"
    CodeNoPermission   = "403"
    CodeNotFound       = "404"
    CodeExceededQPM    = "429"
    CodeServerError    = "500"
)

type BasicError struct {
    Code    string
    Message string
}

func (e *BasicError) Error() string {
    return fmt.Sprintf("error: code=%s, message=%s", e.Code, e.Message)
}

func NewBasicError(code, message string) *BasicError {
    return &BasicError{
        Code:    code,
        Message: message,
    }
}