package usecase

import (
	"context"

	"github.com/phamduytien1805/auth/domain"
)

type LoginUsecase struct {
	userSvc  domain.UserService
	tokenSvc domain.TokenService
}

func NewLoginUsecase(userSvc domain.UserService, tokenSvc domain.TokenService) *LoginUsecase {
	return &LoginUsecase{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
	}
}

func (l *LoginUsecase) Exec(ctx context.Context, identity string, password string) (*domain.User, *domain.TokenPair, error) {
	user, err := l.userSvc.VerifyUserByIdentity(ctx, identity, password)
	if err != nil {
		return nil, nil, err
	}
	tokenPair, err := l.tokenSvc.CreateTokenPair(ctx, user.ID, user.Username, user.Email)
	return user, tokenPair, err

}
