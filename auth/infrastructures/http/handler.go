package http_adapter

import (
	"errors"
	"net/http"

	"github.com/phamduytien1805/auth/domain"
	"github.com/phamduytien1805/auth/interfaces"
	"github.com/phamduytien1805/package/common"
	"github.com/phamduytien1805/package/http_utils"
)

func (s *HttpServer) registerUser(w http.ResponseWriter, r *http.Request) {
	var req interfaces.CreateUserForm
	if !http_utils.DecodeAndValidateRequest(w, r, s.validator, &req) {
		return
	}

	createdUser, tokenPair, err := s.uc.Register.Exec(r.Context(), req.Username, req.Email, req.Credential)
	if err != nil {
		if errors.Is(err, common.ErrorUserResourceConflict) {
			http_utils.EditConflictResponse(w, r, err)
			return
		}
		http_utils.ServerErrorResponse(w, r, err)
		return
	}
	setRfTokenCookie(w, tokenPair.RefreshToken)

	http_utils.Ok(w, r, http.StatusCreated, &interfaces.UserSession{
		User:        createdUser,
		AccessToken: tokenPair.AccessToken,
	})
}

func (s *HttpServer) authenticateUserBasic(w http.ResponseWriter, r *http.Request) {
	var req interfaces.BasicAuthForm
	if !http_utils.DecodeAndValidateRequest(w, r, s.validator, &req) {
		return
	}

	identity := req.Username
	if identity == "" {
		identity = req.Email
	}

	authUser, tokenPair, err := s.uc.Login.Exec(r.Context(), identity, req.Credential)
	if err != nil {
		s.logger.Error("Login failed", "detail", err.Error())
		if errors.Is(err, common.ErrUserNotFound) {
			http_utils.UnauthorizedResponse(w, r, err)
			return
		}
		http_utils.ServerErrorResponse(w, r, err)
		return
	}
	setRfTokenCookie(w, tokenPair.RefreshToken)

	http_utils.Ok(w, r, http.StatusOK, &interfaces.UserSession{
		User:        authUser,
		AccessToken: tokenPair.AccessToken,
	})
}

func (s *HttpServer) refreshToken(w http.ResponseWriter, r *http.Request) {
	rfToken, err := getRfTokenFromCookie(r)
	if err != nil {
		http_utils.UnauthorizedResponse(w, r, err)
		return
	}

	user, tokenPair, err := s.uc.RefreshToken.Exec(r.Context(), rfToken)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}
	setRfTokenCookie(w, tokenPair.RefreshToken)

	http_utils.Ok(w, r, http.StatusOK, &interfaces.UserSession{
		User:        user,
		AccessToken: tokenPair.AccessToken,
	})
}

func (s *HttpServer) logout(w http.ResponseWriter, r *http.Request) {
	rfToken, err := getRfTokenFromCookie(r)
	if err != nil {
		s.logger.Error("No cookies", "detail", err.Error())
		http_utils.BadRequestResponse(w, r, err)
		return

	}
	if rfToken == "" {
		http_utils.Ok(w, r, http.StatusOK, true)
		return
	}

	err = s.uc.Logout.Exec(r.Context(), rfToken)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	http_utils.Ok(w, r, http.StatusOK, true)
}

func (s *HttpServer) verifyEmailUser(w http.ResponseWriter, r *http.Request) {
	var req interfaces.EmailVerificationForm
	if !http_utils.DecodeAndValidateRequest(w, r, s.validator, &req) {
		return
	}

	_, err := s.uc.VerifyEmail.Exec(r.Context(), req.Token)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	http_utils.Ok(w, r, http.StatusOK, true)

}

func (s *HttpServer) resendEmailVerification(w http.ResponseWriter, r *http.Request) {
	tokenPayload := r.Context().Value(authorizationPayloadKey).(domain.TokenPayload)

	err := s.uc.ResendEmail.Exec(r.Context(), tokenPayload.UserID)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	http_utils.Ok(w, r, http.StatusOK, true)
}

func (s *HttpServer) verifyAuthentication(w http.ResponseWriter, r *http.Request) {
	tokenPayload := r.Context().Value(authorizationPayloadKey).(domain.TokenPayload)

	w.Header().Set("X-Forwarded-Userid", tokenPayload.UserID.String())
	http_utils.Ok(w, r, http.StatusOK, true)
}
