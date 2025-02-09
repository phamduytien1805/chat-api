package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/phamduytien1805/hub/interfaces"
	"github.com/phamduytien1805/package/http_utils"
)

func (s *HttpServer) createDmChannel(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(http_utils.UserIdKey).(uuid.UUID)
	var req interfaces.CreateDMForm
	if !http_utils.DecodeAndValidateRequest(w, r, s.validator, &req) {
		return
	}
	dm, err := s.uc.CreateDMChannel.Exec(r.Context(), userId, req.Recipient)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}
	http_utils.Ok(w, r, http.StatusCreated, dm)
}

func (s *HttpServer) getChannelDetail(w http.ResponseWriter, r *http.Request) {
	channelIdParam := chi.URLParam(r, "id")
	channelId, err := uuid.Parse(channelIdParam)
	if err != nil {
		http_utils.BadRequestResponse(w, r, err)
		return
	}
	dm, err := s.uc.GetDMChannel.ById(r.Context(), channelId)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}
	http_utils.Ok(w, r, http.StatusCreated, dm)
}
