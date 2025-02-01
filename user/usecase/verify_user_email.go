package usecase

import (
	"context"

	"github.com/phamduytien1805/user/domain"
)

type VerifyUserEmailUsecase struct {
	repo domain.UserRepo
}

func NewVerifyUserEmailUsecase(userRepo domain.UserRepo) *VerifyUserEmailUsecase {
	return &VerifyUserEmailUsecase{
		repo: userRepo,
	}
}

func (s *VerifyUserEmailUsecase) Exec(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	user.EmailVerified = true
	updatedUser, err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil

}
