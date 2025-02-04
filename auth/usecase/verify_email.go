package usecase

import (
	"context"
	"log/slog"

	"github.com/phamduytien1805/auth/domain"
)

type VerifyEmailUsecase struct {
	mailSvc domain.MailService
	logger  *slog.Logger
}

func NewVerifyEmailUsecase(logger *slog.Logger, mailSvc domain.MailService) *VerifyEmailUsecase {
	return &VerifyEmailUsecase{
		mailSvc: mailSvc,
		logger:  logger,
	}
}

func (uc *VerifyEmailUsecase) Exec(ctx context.Context, token string) (string, error) {
	return uc.mailSvc.VerifyEmail(ctx, token)
}
