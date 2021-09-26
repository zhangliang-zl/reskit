package httpserver

import (
	"net/http"
)

type ErrorInterface interface {
	Error() string
	GetCode() int
}

type HttpError struct {
	code    int
	message string
}

func (h HttpError) Error() string {
	return h.message
}

func (h HttpError) GetCode() int {
	return h.code
}

func NewError(code int, errMsg string) ErrorInterface {
	return HttpError{
		code:    code,
		message: errMsg,
	}
}

func NewBadRequest(errMsg string) ErrorInterface {
	return HttpError{
		code:    http.StatusBadRequest,
		message: errMsg,
	}
}

func NewMethodNotAllowed(errMsg string) ErrorInterface {
	return HttpError{
		code:    http.StatusMethodNotAllowed,
		message: errMsg,
	}
}

func NewInternalError(errMsg string) ErrorInterface {
	return HttpError{
		code:    http.StatusInternalServerError,
		message: errMsg,
	}
}

func NewNotFound(errMsg string) ErrorInterface {
	return HttpError{
		code:    http.StatusNotFound,
		message: errMsg,
	}
}
