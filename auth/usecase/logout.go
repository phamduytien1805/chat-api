package usecase

import (
	"context"
	"log/slog"

	"github.com/phamduytien1805/auth/domain"
)

type LogoutUsecase struct {
	tokenSvc domain.TokenService
	logger   *slog.Logger
}

func NewLogoutUsecase(logger *slog.Logger, tokenSvc domain.TokenService) *LogoutUsecase {
	return &LogoutUsecase{
		tokenSvc: tokenSvc,
		logger:   logger,
	}
}

func (uc *LogoutUsecase) Exec(ctx context.Context, token string) error {
	_, err := uc.tokenSvc.RevokeUserRefreshToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}
