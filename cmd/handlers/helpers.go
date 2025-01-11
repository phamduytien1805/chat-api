package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/phamduytien1805/package/http_utils"
	"github.com/phamduytien1805/package/token"
)

type contextKey string

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = contextKey("authorization_payload")
	authorizationRefreshKey = "refresh_token"
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

		payload, err := s.authSvc.VerifyAccessToken(r.Context(), accessToken)
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

func parseAuthorizationHeader(header string) (authType, token string, err error) {
	fields := strings.Fields(header)
	if len(fields) != 2 {
		return "", "", errors.New("invalid authorization header format")
	}
	return fields[0], fields[1], nil
}
