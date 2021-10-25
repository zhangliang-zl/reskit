package httperror

import (
	"net/http"
)

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) GetCode() int {
	return e.Code
}

func New(code int, msg string) Error {
	return Error{
		Code:    code,
		Message: msg,
	}
}

const (
	MsgNotFound         = "404 not found"
	MsgInternalError    = "500 web internal error"
	MsgMethodNotAllowed = "405 method not allowed"
	MsgForbidden        = "403 forbidden"
	MsgBadRequest       = "400 bad request "
	MsgUnAuthorized     = "401 not authorized "
)

func NewBadRequest(msg ...string) Error {
	message := MsgBadRequest
	if len(msg) > 0 {
		message = msg[0]
	}
	return Error{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewUnauthorized(msg ...string) Error {
	message := MsgUnAuthorized
	if len(msg) > 0 {
		message = msg[0]
	}

	return Error{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewForbidden(msg ...string) Error {
	message := MsgForbidden
	if len(msg) > 0 {
		message = msg[0]
	}

	return Error{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

func NewMethodNotAllowed(msg ...string) Error {
	message := MsgMethodNotAllowed
	if len(msg) > 0 {
		message = msg[0]
	}

	return Error{
		Code:    http.StatusMethodNotAllowed,
		Message: message,
	}
}

func NewInternalError(msg ...string) Error {
	message := MsgInternalError
	if len(msg) > 0 {
		message = msg[0]
	}

	return Error{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func NewNotFound(msg ...string) Error {
	message := MsgNotFound
	if len(msg) > 0 {
		message = msg[0]
	}

	return Error{
		Code:    http.StatusNotFound,
		Message: message,
	}
}
