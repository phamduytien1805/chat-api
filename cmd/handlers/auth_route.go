package handlers

import (
	"errors"
	"net/http"

	"github.com/phamduytien1805/internal/user"
	"github.com/phamduytien1805/package/http_utils"
)

type UserSession struct {
	User        *user.User `json:"user"`
	AccessToken string     `json:"access_token"`
}

// Register a new user
func (s *HttpServer) registerUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserForm
	if !s.decodeAndValidateRequest(w, r, &req) {
		return
	}

	createdUser, err := s.userSvc.CreateUserWithCredential(r.Context(), req)
	if err != nil {
		s.logger.Error(err.Error())
		if errors.Is(err, user.ErrorUserResourceConflict) {
			http_utils.EditConflictResponse(w, r, err)
			return
		}
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	s.createAndSendTokens(w, r, http.StatusCreated, createdUser)
}

// Authenticate a user using basic credentials
func (s *HttpServer) authenticateUserBasic(w http.ResponseWriter, r *http.Request) {
	var req user.BasicAuthForm
	if !s.decodeAndValidateRequest(w, r, &req) {
		return
	}

	authUser, err := s.userSvc.AuthenticateUserBasic(r.Context(), req)
	if err != nil {
		if errors.Is(err, user.ErrorUserInvalidAuthenticate) {
			http_utils.UnauthorizedResponse(w, r, err)
			return
		}
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	s.createAndSendTokens(w, r, http.StatusOK, authUser)
}

// Refresh the access token using the refresh token
func (s *HttpServer) refreshToken(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie(authorizationRefreshKey)
	if err != nil {
		http_utils.BadRequestResponse(w, r, errors.New("refresh token is not provided"))
		return
	}

	payload, err := s.authSvc.VerifyRefreshToken(r.Context(), refreshCookie.Value)
	if err != nil {
		http_utils.BadRequestResponse(w, r, err)
		return
	}

	authUser, err := s.userSvc.GetUserById(r.Context(), payload.UserID)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	s.createAndSendTokens(w, r, http.StatusOK, authUser)
}

func (s *HttpServer) createAndSendTokens(w http.ResponseWriter, r *http.Request, statusCode int, user *user.User) {
	accessToken, err := s.authSvc.CreateAccessTokens(r.Context(), user.ID, user.Username, user.Email)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}
	refreshToken, err := s.authSvc.CreateRefreshTokens(r.Context(), user.ID, user.Username, user.Email)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     authorizationRefreshKey,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
	})

	http_utils.Ok(w, r, statusCode, &UserSession{
		User:        user,
		AccessToken: accessToken,
	})
}
