package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/phamduytien1805/internal/platform/db"
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/hash_generator"
	"github.com/phamduytien1805/package/token"
)

type UserSvc interface {
	CreateUserWithCredential(ctx context.Context, form CreateUserForm) (*User, error)
	AuthenticateUserBasic(ctx context.Context, form BasicAuthForm) (*UserSession, error)
}

type UserSvcImpl struct {
	logger     *slog.Logger
	hashGen    *hash_generator.Argon2idHash
	tokenMaker token.Maker
	config     *config.Config
	repo       userRepo
}

func NewUserServiceImpl(store db.Store, tokenMaker token.Maker, config *config.Config, logger *slog.Logger, hashGen *hash_generator.Argon2idHash) UserSvc {
	return &UserSvcImpl{
		logger:     logger,
		hashGen:    hashGen,
		tokenMaker: tokenMaker,
		config:     config,
		repo:       newUserGatewayImpl(store),
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
		HashedPassword: hashSaltCredential.Hash,
		Salt:           hashSaltCredential.Salt,
	}, func(createdUser *User) error {
		//TODO: add logic to send email verification
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdUser, nil

}

func (s *UserSvcImpl) AuthenticateUserBasic(ctx context.Context, form BasicAuthForm) (*UserSession, error) {
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

	if err = s.hashGen.Compare(userCredential.HashedPassword, userCredential.Salt, form.Credential); err != nil {
		return nil, ErrorUserInvalidAuthenticate
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.Token.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.Token.RefreshTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	return &UserSession{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  *user,
	}, nil
}
