package usecase

import (
	"context"

	"github.com/phamduytien1805/auth/domain"
)

type LogoutUsecase struct {
	tokenSvc domain.TokenService
}

func NewLogoutUsecase(tokenSvc domain.TokenService) *LogoutUsecase {
	return &LogoutUsecase{
		tokenSvc: tokenSvc,
	}
}

func (uc *LogoutUsecase) Exec(ctx context.Context, token string) error {
	_, err := uc.tokenSvc.RevokeUserRefreshToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}
