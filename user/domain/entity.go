package domain

import "github.com/google/uuid"

type User struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
}

type UserCredential struct {
	HashedPassword string
}
