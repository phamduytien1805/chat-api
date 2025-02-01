package usecase

import (
	"context"

	"github.com/phamduytien1805/auth/domain"
)

type VerifyEmailUsecase struct {
	mailSvc domain.MailService
}

func NewVerifyEmailUsecase(mailSvc domain.MailService) *VerifyEmailUsecase {
	return &VerifyEmailUsecase{
		mailSvc: mailSvc,
	}
}

func (uc *VerifyEmailUsecase) Exec(ctx context.Context, token string) (string, error) {
	return uc.mailSvc.VerifyEmail(ctx, token)
}
