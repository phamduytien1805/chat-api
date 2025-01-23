package handlers

import (
	"errors"
	"net/http"

	"github.com/phamduytien1805/internal/auth"
	"github.com/phamduytien1805/internal/user"
	"github.com/phamduytien1805/package/http_utils"
	"github.com/phamduytien1805/package/token"
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

	s.authSvc.SendEmailAsync(r.Context(), createdUser.Email)

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

// Logout current session
func (s *HttpServer) logout(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie(authorizationRefreshKey)
	if err != nil {
		s.logger.Error("No cookies", "detail", err.Error())
	}
	if refreshCookie.Value == "" {
		http_utils.Ok(w, r, http.StatusOK, true)
		return
	}

	_, err = s.authSvc.RevokeUserRefreshToken(r.Context(), refreshCookie.Value)
	if err != nil {
		if !errors.Is(err, auth.ErrRevokedRefreshToken) {
			http_utils.BadRequestResponse(w, r, err)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     authorizationRefreshKey,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
	})

	http_utils.Ok(w, r, http.StatusOK, true)
}

// Refresh the access token using the refresh token
func (s *HttpServer) refreshToken(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie(authorizationRefreshKey)
	if err != nil {
		http_utils.BadRequestResponse(w, r, errors.New("refresh token is not provided"))
		return
	}

	revokedToken, err := s.authSvc.RevokeUserRefreshToken(r.Context(), refreshCookie.Value)
	if err != nil {
		http_utils.BadRequestResponse(w, r, err)
		return
	}

	authUser, err := s.userSvc.GetUserById(r.Context(), revokedToken.UserID)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	s.createAndSendTokens(w, r, http.StatusOK, authUser)
}

// Verify the email using the token
func (s *HttpServer) verifyEmail(w http.ResponseWriter, r *http.Request) {
	var req auth.EmailVerificationForm
	if !s.decodeAndValidateRequest(w, r, &req) {
		return
	}

	_, err := s.authSvc.VerifyEmail(r.Context(), req.Token)
	if err != nil {
		s.logger.Error(err.Error())
		http_utils.BadRequestResponse(w, r, err)
		return
	}

	http_utils.Ok(w, r, http.StatusOK, true)

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

func (s *HttpServer) resendEmailVerification(w http.ResponseWriter, r *http.Request) {
	tokenPayload := r.Context().Value(authorizationPayloadKey).(token.Payload)

	user, err := s.userSvc.GetUserById(r.Context(), tokenPayload.UserID)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	if user.EmailVerified {
		http_utils.EmailVerifiedResponse(w, r, errors.New("email already verified"))
		return
	}

	err = s.authSvc.SendEmailAsync(r.Context(), user.Email)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	http_utils.Ok(w, r, http.StatusOK, true)
}
