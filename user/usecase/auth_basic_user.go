package usecase

import (
	"context"

	"github.com/phamduytien1805/user/domain"
)

type AuthBasicUserUsecase struct {
	repo domain.UserRepo
	hash domain.Hash
}

func NewAuthBasicUserUsecase(userRepo domain.UserRepo, hash domain.Hash) *AuthBasicUserUsecase {
	return &AuthBasicUserUsecase{
		repo: userRepo,
		hash: hash,
	}
}

func (s *AuthBasicUserUsecase) Exec(ctx context.Context, email, credential string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrorUserInvalidAuthenticate
	}

	userCredential, err := s.repo.GetUserCredentialByUserId(ctx, user.ID)
	if err != nil {
		return nil, domain.ErrorUserInvalidAuthenticate
	}

	if err = s.hash.Compare(userCredential.HashedPassword, credential); err != nil {
		return nil, domain.ErrorUserInvalidAuthenticate
	}

	return user, nil

}
