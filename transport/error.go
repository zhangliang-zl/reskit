package transport

import (
	"net/http"
)

type ErrorInterface interface {
	Error() string
	GetCode() int
}

type Error struct {
	Info ErrorInfo `json:"error"`
}

type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
	SubCode int    `json:"sub_code"`
}

func (e Error) Error() string {
	return e.Info.Message
}

func (e Error) GetCode() int {
	return e.Info.Code
}

func NewError(code int, errMsg string, subCode int) Error {
	return Error{
		Info: ErrorInfo{
			Code:    code,
			Message: errMsg,
			Type:    "",
			SubCode: subCode,
		}}
}

func NewInternalError(errMsg string) Error {
	return NewError(http.StatusInternalServerError, errMsg, http.StatusInternalServerError)
}
