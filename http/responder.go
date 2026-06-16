package http

import (
	gohttp "net/http"

	"github.com/dani-susanto/go-common/json"
)

type Responder interface {
	Success(w gohttp.ResponseWriter, statusCode int, data any, message string, meta any)
	Error(w gohttp.ResponseWriter, statusCode int, message string, errors any)
}

type responder struct {
	json json.JSON
}

func NewResponder(json json.JSON) Responder {
	return &responder{
		json: json,
	}
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func (r *responder) Success(w gohttp.ResponseWriter, statusCode int, data any, message string, meta any) {
	r.output(w, statusCode, true, data, message, meta, nil)
}

func (r *responder) Error(w gohttp.ResponseWriter, statusCode int, message string, errors any) {
	r.output(w, statusCode, false, nil, message, nil, errors)
}

// func (r *responder) output(w gohttp.ResponseWriter, statusCode int, status bool, data any, message string, meta any, errors any) {
// 	res := response{
// 		Status:  status,
// 		Message: message,
// 		Data:    data,
// 		Meta:    meta,
// 		Errors:  errors,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	if err := r.json.EncodeHttp(w, res); err != nil {
// 		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
// 		return
// 	}
// }

func (r *responder) output(w gohttp.ResponseWriter, statusCode int, status bool, data any, message string, meta any, errors any) {
	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
		Meta:    meta,
		Errors:  errors,
	}

	body, err := r.json.Encode(res)
	if err != nil {
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}
