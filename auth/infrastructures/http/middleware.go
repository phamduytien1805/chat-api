package http_adapter

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/phamduytien1805/package/http_utils"
	"github.com/phamduytien1805/package/token"
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

		payload, err := s.uc.VerifyAccessToken.Exec(accessToken)
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
