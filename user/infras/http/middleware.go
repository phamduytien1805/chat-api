package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	userIdKey = contextKey("userId")
)

// Authenticator middleware to validate and attach authorization payload to context
func (s *HttpServer) extractForwardedHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := uuid.Parse(r.Header.Get("X-Forwarded-Userid"))
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, userIdKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
