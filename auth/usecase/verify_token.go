package usecase

import (
	"log/slog"

	"github.com/phamduytien1805/auth/domain"
)

type VerifyAccessTokenUsecase struct {
	tokenSvc domain.TokenService
	logger   *slog.Logger
}

func NewVerifyAccessTokenUsecase(logger *slog.Logger, tokenSvc domain.TokenService) *VerifyAccessTokenUsecase {
	return &VerifyAccessTokenUsecase{
		tokenSvc: tokenSvc,
		logger:   logger,
	}
}

func (uc *VerifyAccessTokenUsecase) Exec(token string) (*domain.TokenPayload, error) {
	user, err := uc.tokenSvc.VerifyAccessToken(token)
	return user, err
}
