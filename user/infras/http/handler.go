package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/phamduytien1805/package/http_utils"
)

func (s *HttpServer) getUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(http_utils.UserIdKey).(uuid.UUID)

	user, err := s.uc.GetUser.ById(r.Context(), userId)
	if err != nil {
		http_utils.BadRequestResponse(w, r, err)
		return
	}
	http_utils.Ok(w, r, http.StatusOK, user)
}
