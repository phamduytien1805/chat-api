package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/phamduytien1805/package/token"
)

type TokenPayload struct {
	*token.Payload
}
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenService interface {
	VerifyAccessToken(token string) (*TokenPayload, error)
	RevokeUserRefreshToken(ctx context.Context, tokenString string) (*TokenPayload, error)
	InvalidateRefreshTokenIfNeeded(ctx context.Context, tokenString string, tokenExpiredAt time.Time) error
	CreateTokenPair(ctx context.Context, userID uuid.UUID, username string, email string) (*TokenPair, error)
}

var (
	ErrRevokedRefreshToken = errors.New("refresh token is already used")
)
