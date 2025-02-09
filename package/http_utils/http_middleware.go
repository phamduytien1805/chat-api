package http_utils

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIdKey = contextKey("userId")
)

// Authenticator middleware to validate and attach authorization payload to context
func ExtractForwardedHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := uuid.Parse(r.Header.Get("X-Forwarded-Userid"))
		if err != nil {
			BadRequestResponse(w, r, errors.New("invalid user id"))
			return
		}
		ctx = context.WithValue(ctx, UserIdKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
