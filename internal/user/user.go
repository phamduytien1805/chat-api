package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/phamduytien1805/internal/platform/db"
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/hash_generator"
)

type UserSvc interface {
	CreateUserWithCredential(ctx context.Context, form CreateUserForm) (*User, error)
	AuthenticateUserBasic(ctx context.Context, form BasicAuthForm) (*User, error)
	GetUserById(ctx context.Context, userID uuid.UUID) (*User, error)
}

type UserSvcImpl struct {
	logger  *slog.Logger
	hashGen *hash_generator.Argon2idHash
	config  *config.Config
	repo    userRepo
}

func NewUserServiceImpl(store db.Store, config *config.Config, logger *slog.Logger, hashGen *hash_generator.Argon2idHash) UserSvc {
	return &UserSvcImpl{
		logger:  logger,
		hashGen: hashGen,
		config:  config,
		repo:    newUserGatewayImpl(store),
	}
}

func (s *UserSvcImpl) CreateUserWithCredential(ctx context.Context, form CreateUserForm) (*User, error) {
	ID, err := uuid.NewV7()

	if err != nil {
		s.logger.Error("generate uuid user failed", "detail", err.Error())
		return nil, err
	}

	hashSaltCredential, err := s.hashGen.GenerateHash([]byte(form.Credential), nil)
	if err != nil {
		s.logger.Error("error while hashing password user", "detail", err.Error())
		return nil, err
	}

	createdUser, err := s.repo.createUserWithCredential(ctx, &User{
		ID:            ID,
		Username:      form.Username,
		Email:         form.Email,
		EmailVerified: false,
	}, &UserCredential{
		HashedPassword: hashSaltCredential,
	}, func(createdUser *User) error {
		//TODO: add logic to send email verification
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserSvcImpl) AuthenticateUserBasic(ctx context.Context, form BasicAuthForm) (*User, error) {
	user, err := s.repo.getUserByEmail(ctx, form.Email)
	if err != nil {
		s.logger.Error("error getting user by email", "detail", err.Error())
		return nil, ErrorUserInvalidAuthenticate
	}

	userCredential, err := s.repo.getUserCredentialByUserId(ctx, user.ID)
	if err != nil {
		s.logger.Error("error getting user credential", "detail", err.Error())
		return nil, ErrorUserInvalidAuthenticate
	}

	if err = s.hashGen.Compare(userCredential.HashedPassword, form.Credential); err != nil {
		return nil, ErrorUserInvalidAuthenticate
	}

	return user, nil
}

func (s *UserSvcImpl) GetUserById(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := s.repo.getUserById(ctx, userID)
	if err != nil {
		s.logger.Error("error getting user by id", "detail", err.Error())
		return nil, err
	}

	return user, nil
}
