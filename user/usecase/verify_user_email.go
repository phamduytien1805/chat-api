package usecase

import (
	"context"
	"log/slog"

	"github.com/phamduytien1805/user/domain"
)

type VerifyUserEmailUsecase struct {
	repo   domain.UserRepo
	logger *slog.Logger
}

func NewVerifyUserEmailUsecase(logger *slog.Logger, userRepo domain.UserRepo) *VerifyUserEmailUsecase {
	return &VerifyUserEmailUsecase{
		repo:   userRepo,
		logger: logger,
	}
}

func (s *VerifyUserEmailUsecase) Exec(ctx context.Context, email string) (domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	user.EmailVerified = true
	updatedUser, err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return updatedUser, nil

}
