package interfaces

import (
	"github.com/phamduytien1805/auth/domain"
)

type BasicAuthForm struct {
	Username   string `json:"username" validate:"required_without=Email,omitempty,min=5,max=32"`
	Email      string `json:"email" validate:"required_without=Username,omitempty,email"`
	Credential string `json:"password" validate:"required,min=8"`
}

type CreateUserForm struct {
	Username   string `json:"username" validate:"required,min=2,max=32"`
	Email      string `json:"email" validate:"required,email"`
	Credential string `json:"password" validate:"required,min=8"`
}

type EmailVerificationForm struct {
	Token string `json:"token" validate:"required"`
}

type UserSession struct {
	User        *domain.User `json:"user"`
	AccessToken string       `json:"access_token"`
}
