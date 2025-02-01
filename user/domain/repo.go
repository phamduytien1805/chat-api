package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepo interface {
	CreateUserWithCredential(ctx context.Context, userParams *User, userCredential *UserCredential) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*User, error)
	GetUserById(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserCredentialByUserId(ctx context.Context, userID uuid.UUID) (*UserCredential, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
}
