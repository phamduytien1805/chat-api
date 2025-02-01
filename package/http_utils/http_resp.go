package http_utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/phamduytien1805/package/validator"
)

type envelope map[string]any

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func Ok(w http.ResponseWriter, r *http.Request, status int, payload any) {
	err := writeJSON(w, status, envelope{"data": payload, "code": 0}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func DecodeAndValidateRequest(w http.ResponseWriter, r *http.Request, validator *validator.Validate, req interface{}) bool {
	if err := render.DecodeJSON(r.Body, req); err != nil {
		BadRequestResponse(w, r, err)
		return false
	}

	if err := validator.Struct(req); err != nil {
		FailedValidationResponse(w, r, err)
		return false
	}

	return true
}
