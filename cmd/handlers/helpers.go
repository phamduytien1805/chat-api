package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/phamduytien1805/package/http_utils"
)

func (s *HttpServer) decodeAndValidateRequest(w http.ResponseWriter, r *http.Request, req interface{}) bool {
	if err := render.DecodeJSON(r.Body, req); err != nil {
		http_utils.BadRequestResponse(w, r, err)
		return false
	}

	if err := s.validator.Struct(req); err != nil {
		http_utils.FailedValidationResponse(w, r, err)
		return false
	}

	return true
}
