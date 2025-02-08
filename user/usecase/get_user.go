package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/phamduytien1805/package/common"
	"github.com/phamduytien1805/user/domain"
)

type GetUserUsecase struct {
	repo   domain.UserRepo
	logger *slog.Logger
}

func NewGetUserUsecase(logger *slog.Logger, userRepo domain.UserRepo) *GetUserUsecase {
	return &GetUserUsecase{
		repo:   userRepo,
		logger: logger,
	}
}

func (s *GetUserUsecase) ById(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	user, err := s.repo.GetUserById(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s *GetUserUsecase) ByEmailOrUsername(ctx context.Context, emailOrUsername string) (domain.User, error) {
	user, err := s.repo.GetUserByEmailOrUsername(ctx, emailOrUsername)
	if err != nil {
		return domain.User{}, common.ErrUserNotFound
	}

	return user, nil
}
