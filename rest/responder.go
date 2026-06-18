package rest

import (
	"fmt"
	"net/http"

	"github.com/dani-susanto/go-common/json"
)

type Responder interface {
	Success(w http.ResponseWriter, statusCode int, data any, message string, meta *MetaData)
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

type MetaData struct {
	CurrentPage int     `json:"current_page"`
	PerPage     int     `json:"per_page"`
	Total       int     `json:"total"`
	LastPage    int     `json:"last_page"`
	From        int     `json:"from"`
	To          int     `json:"to"`
	NextPage    *int    `json:"next_page"`
	PrevPage    *int    `json:"prev_page"`
	NextLink    *string `json:"next_link"`
	PrevLink    *string `json:"prev_link"`
}

type SuccessResponse struct {
	Status  bool      `json:"status" example:"true"`
	Message string    `json:"message,omitempty"`
	Data    any       `json:"data,omitempty"`
	Meta    *MetaData `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Status  bool                `json:"status" example:"false"`
	Message string              `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

func (r *responder) Success(w http.ResponseWriter, statusCode int, data any, message string, meta *MetaData) {
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
		Errors:  normalizeErrors(errors),
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

func normalizeErrors(errors any) map[string][]string {
	if errors == nil {
		return nil
	}
	switch v := errors.(type) {
	case map[string][]string:
		return v
	case map[string]string:
		result := make(map[string][]string)
		for k, val := range v {
			result[k] = []string{val}
		}
		return result
	case string:
		return map[string][]string{"message": {v}}
	case []string:
		return map[string][]string{"message": v}
	case error:
		return map[string][]string{"message": {v.Error()}}
	default:
		return map[string][]string{"message": {"unknown error"}}
	}
}

func BuildMetaData(page, perPage, total int, baseURL string) *MetaData {
	lastPage := total / perPage
	if total%perPage > 0 {
		lastPage++
	}
	from := (page-1)*perPage + 1
	to := from + perPage - 1
	if to > total {
		to = total
	}

	var nextPage *int
	var prevPage *int
	var nextLink *string
	var prevLink *string

	if page < lastPage {
		next := page + 1
		nextPage = &next
		link := fmt.Sprintf("%s?page=%d&per_page=%d", baseURL, next, perPage)
		nextLink = &link
	}

	if page > 1 {
		prev := page - 1
		prevPage = &prev
		link := fmt.Sprintf("%s?page=%d&per_page=%d", baseURL, prev, perPage)
		prevLink = &link
	}

	return &MetaData{
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		LastPage:    lastPage,
		From:        from,
		To:          to,
		NextPage:    nextPage,
		PrevPage:    prevPage,
		NextLink:    nextLink,
		PrevLink:    prevLink,
	}
}
