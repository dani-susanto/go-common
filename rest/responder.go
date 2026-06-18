package rest

import (
	"net/http"

	"github.com/dani-susanto/go-common/json"
)

type Responder interface {
	Success(w http.ResponseWriter, statusCode int, data any, message string, meta any)
	Error(w http.ResponseWriter, statusCode int, message string, errors any)
}

type responder struct {
	json json.JSON
}

func NewResponder(json json.JSON) Responder {
	return &responder{
		json: json,
	}
}

type SuccessResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message,omitempty" example:"login successful"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Status  bool   `json:"status" example:"false"`
	Message string `json:"message,omitempty" example:"user not found"`
	Errors  any    `json:"errors,omitempty"`
}

func (r *responder) Success(w http.ResponseWriter, statusCode int, data any, message string, meta any) {
	r.write(w, statusCode, SuccessResponse{
		Status:  true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func (r *responder) Error(w http.ResponseWriter, statusCode int, message string, errors any) {
	r.write(w, statusCode, ErrorResponse{
		Status:  false,
		Message: message,
		Errors:  errors,
	})
}

func (r *responder) write(w http.ResponseWriter, statusCode int, res any) {
	body, err := r.json.Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}
