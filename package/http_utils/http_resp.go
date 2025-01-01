package http_utils

import (
	"encoding/json"
	"net/http"
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
