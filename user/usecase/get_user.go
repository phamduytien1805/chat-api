package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/phamduytien1805/user/domain"
)

type GetUserUsecase struct {
	repo domain.UserRepo
}

func NewGetUserUsecase(userRepo domain.UserRepo) *GetUserUsecase {
	return &GetUserUsecase{
		repo: userRepo,
	}
}

func (s *GetUserUsecase) ById(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *GetUserUsecase) ByEmailOrUsername(ctx context.Context, emailOrUsername string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmailOrUsername(ctx, emailOrUsername)
	if err != nil {
		return nil, err
	}

	return user, nil
}
