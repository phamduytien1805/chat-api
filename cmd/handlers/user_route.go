package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/phamduytien1805/internal/user"
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

func (s *HttpServer) registerUser(w http.ResponseWriter, r *http.Request) {
	createUserRequest := &user.CreateUserForm{}
	if err := render.DecodeJSON(r.Body, createUserRequest); err != nil {
		s.logger.Error(err.Error())
		http_utils.BadRequestResponse(w, r, err)
		return
	}

	if err := s.validator.Struct(createUserRequest); err != nil {
		s.logger.Error(err.Error())
		http_utils.FailedValidationResponse(w, r, err)
		return
	}

	createdUser, err := s.userSvc.CreateUserWithCredential(r.Context(), *createUserRequest)
	if err != nil {
		s.logger.Error(err.Error())
		if errors.Is(err, user.ErrorUserResourceConflict) {
			http_utils.EditConflictResponse(w, r, err)
			return
		}
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	http_utils.Ok(w, r, http.StatusCreated, createdUser)
}

func (s *HttpServer) authenticateUserBasic(w http.ResponseWriter, r *http.Request) {
	basicAuthForm := &user.BasicAuthForm{}
	if err := render.DecodeJSON(r.Body, basicAuthForm); err != nil {
		http_utils.BadRequestResponse(w, r, err)
		return
	}

	if err := s.validator.Struct(basicAuthForm); err != nil {
		http_utils.FailedValidationResponse(w, r, err)
		return
	}

	userSession, err := s.userSvc.AuthenticateUserBasic(r.Context(), *basicAuthForm)
	if err != nil {
		if errors.Is(err, user.ErrorUserInvalidAuthenticate) {
			http_utils.InvalidAuthenticateResponse(w, r, err)
			return
		}
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	http_utils.Ok(w, r, http.StatusCreated, userSession)

}
