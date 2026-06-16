package json

import (
	gojson "encoding/json"
	"net/http"
)

type JSON interface {
	Encode(v any) (string, error)
	Decode(data []byte, v any) error
	EncodeHttp(w http.ResponseWriter, v any) error
	DecodeHttp(r *http.Request, v any) error
}

type json struct{}

func (j *json) Decode(data []byte, v any) error {
	if err := gojson.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

func (j *json) Encode(v any) (string, error) {
	res, err := gojson.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (j *json) EncodeHttp(w http.ResponseWriter, v any) error {
	if err := gojson.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

func (j *json) DecodeHttp(r *http.Request, v any) error {
	if err := gojson.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func New() JSON {
	return &json{}
}
