package handlers

import (
	"net/http"

	"github.com/phamduytien1805/package/http_utils"
	"github.com/phamduytien1805/package/token"
)

func (s *HttpServer) getUser(w http.ResponseWriter, r *http.Request) {
	tokenPayload := r.Context().Value(authorizationPayloadKey).(token.Payload)
	user, err := s.userSvc.GetUserById(r.Context(), tokenPayload.UserID)
	if err != nil {
		s.logger.Error(err.Error())
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	http_utils.Ok(w, r, http.StatusOK, user)
}
