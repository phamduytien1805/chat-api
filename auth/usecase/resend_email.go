package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/phamduytien1805/auth/domain"
)

type ResendEmailUsecase struct {
	mailSvc domain.MailService
	userSvc domain.UserService
	logger  *slog.Logger
}

func NewResendEmailUsecase(logger *slog.Logger, mailSvc domain.MailService, userSvc domain.UserService) *ResendEmailUsecase {
	return &ResendEmailUsecase{
		mailSvc: mailSvc,
		userSvc: userSvc,
	}
}

func (uc *ResendEmailUsecase) Exec(ctx context.Context, userId uuid.UUID) error {
	user, err := uc.userSvc.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user.EmailVerified {
		return errors.New("email already verified")
	}
	return uc.mailSvc.SendEmailAsync(ctx, user.Email)
}
