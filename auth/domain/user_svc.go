package domain

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
}

type UserService interface {
	CreateUserWithCredential(ctx context.Context, username string, email string, hashed_password string) (*User, error)
	VerifyUserByIdentity(ctx context.Context, identity string, hashed_password string) (*User, error)
	GetUserById(ctx context.Context, userID uuid.UUID) (*User, error)
	VerifyUserEmail(ctx context.Context, userEmail string) error
}
