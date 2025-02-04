package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/phamduytien1805/user/domain"
)

type CreateUserUsecase struct {
	repo   domain.UserRepo
	hash   domain.Hash
	logger *slog.Logger
}

func NewCreateUserUsecase(logger *slog.Logger, userRepo domain.UserRepo, hash domain.Hash) *CreateUserUsecase {
	return &CreateUserUsecase{
		repo:   userRepo,
		hash:   hash,
		logger: logger,
	}
}

func (s *CreateUserUsecase) Exec(ctx context.Context, username, email, credential string) (*domain.User, error) {
	ID, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	hashSaltCredential, err := s.hash.GenerateHash([]byte(credential), nil)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.repo.CreateUserWithCredential(ctx, &domain.User{
		ID:            ID,
		Username:      username,
		Email:         email,
		EmailVerified: false,
	}, &domain.UserCredential{
		HashedPassword: hashSaltCredential,
	})

	if err != nil {
		return nil, err
	}

	return createdUser, nil

}
