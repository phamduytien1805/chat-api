package usecase

import (
	"context"

	"github.com/phamduytien1805/auth/domain"
)

type RegisterUsecase struct {
	userSvc  domain.UserService
	mailSvc  domain.MailService
	tokenSvc domain.TokenService
}

func NewRegisterUsecase(userSvc domain.UserService, mailSvc domain.MailService, tokenSvc domain.TokenService) *RegisterUsecase {
	return &RegisterUsecase{
		userSvc:  userSvc,
		mailSvc:  mailSvc,
		tokenSvc: tokenSvc,
	}
}

func (r *RegisterUsecase) Exec(ctx context.Context, username string, email string, password string) (*domain.User, *domain.TokenPair, error) {
	user, err := r.userSvc.CreateUserWithCredential(ctx, username, email, password)
	if err != nil {
		return nil, nil, err
	}
	tokenPair, err := r.tokenSvc.CreateTokenPair(ctx, user.ID, user.Username, user.Email)
	r.mailSvc.SendEmailAsync(ctx, email)
	return user, tokenPair, err
}
