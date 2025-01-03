// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserCredential(ctx context.Context, arg CreateUserCredentialParams) (UserCredential, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserCredentialByUserId(ctx context.Context, userID uuid.UUID) (UserCredential, error)
}

var _ Querier = (*Queries)(nil)
