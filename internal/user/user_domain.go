package user

import (
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
}

type CreateUserForm struct {
	Username   string `json:"username" validate:"required,min=2,max=32"`
	Email      string `json:"email" validate:"required,email"`
	Credential string `json:"password" validate:"required,min=8"`
}

type BasicAuthForm struct {
	Username   string `json:"username" validate:"required_without=Email,omitempty,min=5,max=32"`
	Email      string `json:"email" validate:"required_without=Username,omitempty,email"`
	Credential string `json:"password" validate:"required,min=8"`
}

type UserCredential struct {
	HashedPassword string
}
