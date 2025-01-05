package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/phamduytien1805/internal/user"
	"github.com/phamduytien1805/package/http_utils"
	"github.com/phamduytien1805/package/token"
)

type UserSession struct {
	User        *user.User `json:"user"`
	AccessToken string     `json:"access_token"`
}

type contextKey string

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = contextKey("authorization_payload")
	authorizationRefreshKey = "refresh_token"
)

// Authenticator middleware to validate and attach authorization payload to context
func (s *HttpServer) authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authorizationHeader := r.Header.Get(authorizationHeaderKey)

		if authorizationHeader == "" {
			http_utils.InvalidAuthenticateResponse(w, r, errors.New("missing authorization header"))
			return
		}

		authorizationType, accessToken, err := parseAuthorizationHeader(authorizationHeader)
		if err != nil || strings.ToLower(authorizationType) != authorizationTypeBearer {
			http_utils.InvalidAuthenticateResponse(w, r, fmt.Errorf("unsupported or invalid authorization type: %s", authorizationType))
			return
		}

		payload, err := s.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			s.logger.Error(err.Error())
			if errors.Is(err, token.ErrInvalidToken) {
				http_utils.TokenExpired(w, r, err)
			} else {
				http_utils.InvalidAuthenticateResponse(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, authorizationPayloadKey, *payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

	http_utils.Ok(w, r, http.StatusCreated, createdUser)
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
			http_utils.InvalidAuthenticateResponse(w, r, err)
			return
		}
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	s.createAndSendTokens(w, r, authUser)
}

// Refresh the access token using the refresh token
func (s *HttpServer) refreshToken(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie(authorizationRefreshKey)
	if err != nil {
		http_utils.InvalidAuthenticateResponse(w, r, errors.New("refresh token is not provided"))
		return
	}

	payload, err := s.tokenMaker.VerifyToken(refreshCookie.Value)
	if err != nil {
		http_utils.InvalidAuthenticateResponse(w, r, err)
		return
	}

	authUser, err := s.userSvc.GetUserById(r.Context(), payload.UserID)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	s.createAndSendTokens(w, r, authUser)
}

func (s *HttpServer) createAndSendTokens(w http.ResponseWriter, r *http.Request, user *user.User) {
	accessToken, refreshToken, err := s.makeToken(user)
	if err != nil {
		http_utils.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     authorizationRefreshKey,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   s.config.Env == "production",
		SameSite: http.SameSiteNoneMode,
	})

	http_utils.Ok(w, r, http.StatusCreated, &UserSession{
		User:        user,
		AccessToken: accessToken,
	})
}

func parseAuthorizationHeader(header string) (authType, token string, err error) {
	fields := strings.Fields(header)
	if len(fields) != 2 {
		return "", "", errors.New("invalid authorization header format")
	}
	return fields[0], fields[1], nil
}

func (s *HttpServer) makeToken(user *user.User) (string, string, error) {
	accessToken, _, err := s.tokenMaker.CreateToken(user.ID, user.Username, user.Email, s.config.Token.AccessTokenDuration)
	if err != nil {
		return "", "", err
	}

	refreshToken, _, err := s.tokenMaker.CreateToken(user.ID, user.Username, user.Email, s.config.Token.RefreshTokenDuration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
