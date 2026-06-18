package exception

import (
	"net/http"
)

type Status string

const (
	INTERNAL_SERVER_ERROR Status = "INTERNAL_SERVER_ERROR"
	BAD_REQUEST           Status = "BAD_REQUEST"
	UNAUTHORIZED          Status = "UNAUTHORIZED"
	FORBIDDEN             Status = "FORBIDDEN"
	NOT_FOUND             Status = "NOT_FOUND"
	CONFLICT              Status = "CONFLICT"
	UNPROCESSABLE_ENTITY  Status = "UNPROCESSABLE_ENTITY"
)

type exception struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
}

func (e *exception) Error() string {
	return e.Message
}

func Throw(status Status, message string) error {
	return &exception{
		Status:  status,
		Message: message,
	}
}

func AsException(err error) (*exception, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(*exception)
	return e, ok
}

var statusToHTTP = map[Status]int{
	INTERNAL_SERVER_ERROR: http.StatusInternalServerError,
	BAD_REQUEST:           http.StatusBadRequest,
	UNAUTHORIZED:          http.StatusUnauthorized,
	FORBIDDEN:             http.StatusForbidden,
	NOT_FOUND:             http.StatusNotFound,
	CONFLICT:              http.StatusConflict,
	UNPROCESSABLE_ENTITY:  http.StatusUnprocessableEntity,
}

func GetHttpCode(err error) int {
	if e, ok := AsException(err); ok {
		if code, ok := statusToHTTP[e.Status]; ok {
			return code
		}
	}
	return http.StatusInternalServerError
}
