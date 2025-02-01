package tokensvc

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phamduytien1805/auth/domain"
	"github.com/phamduytien1805/package/config"
	redis_engine "github.com/phamduytien1805/package/redis"
	"github.com/phamduytien1805/package/token"
)

type TokenService struct {
	tokenMaker           token.Maker
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	redis                redis_engine.RedisQuerier
}

func NewTokenService(config *config.TokenConfig, redis redis_engine.RedisQuerier) (domain.TokenService, error) {
	tokenMaker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		return nil, err
	}
	return &TokenService{
		tokenMaker:           tokenMaker,
		accessTokenDuration:  config.AccessTokenDuration,
		refreshTokenDuration: config.RefreshTokenDuration,
		redis:                redis,
	}, nil
}

func (svc *TokenService) CreateTokenPair(ctx context.Context, userID uuid.UUID, username string, email string) (*domain.TokenPair, error) {
	accessToken, _, err := svc.tokenMaker.CreateToken(userID, username, email, svc.accessTokenDuration)
	if err != nil {
		return nil, err
	}
	refreshToken, _, err := svc.tokenMaker.CreateToken(userID, username, email, svc.refreshTokenDuration)
	if err != nil {
		return nil, err
	}
	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (svc *TokenService) VerifyAccessToken(tokenString string) (*domain.TokenPayload, error) {
	payload, err := svc.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	return &domain.TokenPayload{Payload: payload}, nil
}

func (svc *TokenService) RevokeUserRefreshToken(ctx context.Context, tokenString string) (*domain.TokenPayload, error) {
	payload, err := svc.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	return &domain.TokenPayload{Payload: payload}, nil
}

func (svc *TokenService) InvalidateRefreshTokenIfNeeded(ctx context.Context, tokenString string, tokenExpiredAt time.Time) error {
	hashedToken := sha256.Sum256([]byte(tokenString))
	blacklistKey := fmt.Sprintf("invalid_rftoken:%x", hashedToken)
	keyExist, err := svc.redis.Exist(ctx, blacklistKey)
	if err != nil {
		return err
	}
	if keyExist {
		return domain.ErrRevokedRefreshToken
	}
	// Invalidate the token
	if err := svc.redis.SetTx(ctx, blacklistKey, 1, time.Until(tokenExpiredAt)); err != nil {
		return err
	}
	return nil
}
