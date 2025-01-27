package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/phamduytien1805/internal/platform/db"
	"github.com/phamduytien1805/internal/platform/mail"
	"github.com/phamduytien1805/internal/platform/redis_engine"
	"github.com/phamduytien1805/internal/taskq"
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/token"
)

const (
	emailVerificationKey = "verify_email"
)

type AuthService interface {
	VerifyAccessToken(ctx context.Context, tokenString string) (*token.Payload, error)
	RevokeUserRefreshToken(ctx context.Context, tokenString string) (*token.Payload, error)
	CreateAccessTokens(ctx context.Context, userID uuid.UUID, username string, email string) (string, error)
	CreateRefreshTokens(ctx context.Context, userID uuid.UUID, username string, email string) (string, error)
	SendEmailAsync(ctx context.Context, userEmail string) error
	VerifyEmail(ctx context.Context, token string) (string, error)
}

type AuthServiceImpl struct {
	tokenMaker token.Maker
	redis      redis_engine.RedisQuerier
	logger     *slog.Logger
	taskq      taskq.TaskProducer
	repo       authRepo

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	emailDuration        time.Duration
	verifyLink           string
}

func NewAuthService(config *config.Config, logger *slog.Logger, tokenMaker token.Maker, taskqProducer taskq.TaskProducer, redis redis_engine.RedisQuerier, store db.Store) AuthService {
	return &AuthServiceImpl{
		tokenMaker:           tokenMaker,
		redis:                redis,
		logger:               logger,
		taskq:                taskqProducer,
		repo:                 newAuthGatewayImpl(store),
		accessTokenDuration:  config.Token.AccessTokenDuration,
		refreshTokenDuration: config.Token.RefreshTokenDuration,
		emailDuration:        config.Mail.Expired,
		verifyLink:           config.Web.Http.Server.VerifyEmailUrl,
	}
}

func (a *AuthServiceImpl) VerifyAccessToken(ctx context.Context, tokenString string) (*token.Payload, error) {
	payload, err := a.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (a *AuthServiceImpl) RevokeUserRefreshToken(ctx context.Context, tokenString string) (*token.Payload, error) {
	payload, err := a.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check blacklist for used tokens
	if err := a.invalidateRefreshTokenIfNeeded(ctx, tokenString, payload.ExpiredAt); err != nil {
		return nil, err
	}

	return payload, nil
}

func (a *AuthServiceImpl) CreateAccessTokens(ctx context.Context, userID uuid.UUID, username string, email string) (string, error) {
	accessToken, _, err := a.tokenMaker.CreateToken(userID, username, email, a.accessTokenDuration)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
func (a *AuthServiceImpl) CreateRefreshTokens(ctx context.Context, userID uuid.UUID, username string, email string) (string, error) {
	rfToken, _, err := a.tokenMaker.CreateToken(userID, username, email, a.refreshTokenDuration)
	if err != nil {
		return "", err
	}
	return rfToken, nil
}

func (a *AuthServiceImpl) SendEmailAsync(ctx context.Context, userEmail string) error {
	token, err := uuid.NewRandom()
	if err != nil {
		a.logger.Error("failed to generate uuid SendEmailAsync", "detail", err.Error())
		return err
	}
	verifyKey := fmt.Sprintf("%s:%s", emailVerificationKey, token)
	if err := a.redis.SetTx(ctx, verifyKey, userEmail, a.emailDuration); err != nil {
		a.logger.Error("failed to set verifyKey", "detail", err.Error())
		return err
	}
	payload := mail.SendMailParams{
		To:         userEmail,
		VerifyLink: fmt.Sprintf("%s?token=%s", a.verifyLink, token),
	}
	if err := a.taskq.EnqueueSendMailTask(ctx, payload, asynq.MaxRetry(5)); err != nil {
		return err
	}
	return nil
}

func (a *AuthServiceImpl) VerifyEmail(ctx context.Context, token string) (string, error) {
	verifyKey := fmt.Sprintf("%s:%s", emailVerificationKey, token)
	userEmail, err := a.redis.GetRaw(ctx, verifyKey)
	if err != nil {
		return "", err
	}
	if err := a.redis.Delete(ctx, verifyKey); err != nil {
		return "", err
	}
	if err := a.repo.verifyUserEmail(ctx, userEmail); err != nil {
		return "", err
	}
	return userEmail, nil
}

func (a *AuthServiceImpl) invalidateRefreshTokenIfNeeded(ctx context.Context, tokenString string, tokenExpiredAt time.Time) error {
	hashedToken := sha256.Sum256([]byte(tokenString))
	blacklistKey := fmt.Sprintf("invalid_rftoken:%x", hashedToken)
	keyExist, err := a.redis.Exist(ctx, blacklistKey)
	if err != nil {
		return err
	}
	if keyExist {
		return ErrRevokedRefreshToken
	}
	// Invalidate the token
	if err := a.redis.SetTx(ctx, blacklistKey, 1, time.Until(tokenExpiredAt)); err != nil {
		return err
	}
	return nil
}
