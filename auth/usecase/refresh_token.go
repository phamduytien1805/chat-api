package usecase

import (
	"context"

	"github.com/phamduytien1805/auth/domain"
)

type RefreshTokenUsecase struct {
	tokenSvc domain.TokenService
	userSvc  domain.UserService
}

func NewRefreshTokenUsecase(tokenSvc domain.TokenService, userSvc domain.UserService) *RefreshTokenUsecase {
	return &RefreshTokenUsecase{
		tokenSvc: tokenSvc,
		userSvc:  userSvc,
	}
}

func (uc *RefreshTokenUsecase) Exec(ctx context.Context, rfToken string) (*domain.User, *domain.TokenPair, error) {
	revokedToken, err := uc.tokenSvc.RevokeUserRefreshToken(ctx, rfToken)
	if err != nil {
		return nil, nil, err
	}
	user, err := uc.userSvc.GetUserById(ctx, revokedToken.UserID)
	if err != nil {
		return nil, nil, err
	}
	tokenPair, err := uc.tokenSvc.CreateTokenPair(ctx, user.ID, user.Username, user.Email)
	return user, tokenPair, err

}
