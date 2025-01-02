// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  id,
  username,
  email,
  email_verified
) VALUES (
  $1, $2, $3, $4
) RETURNING id, username, email, email_verified, created_at
`

type CreateUserParams struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.EmailVerified,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.EmailVerified,
		&i.CreatedAt,
	)
	return i, err
}

const createUserCredential = `-- name: CreateUserCredential :one
INSERT INTO user_credentials (
  user_id,
  hashed_password
) VALUES (
  $1, $2
) RETURNING user_id, hashed_password
`

type CreateUserCredentialParams struct {
	UserID         uuid.UUID `json:"user_id"`
	HashedPassword string    `json:"hashed_password"`
}

func (q *Queries) CreateUserCredential(ctx context.Context, arg CreateUserCredentialParams) (UserCredential, error) {
	row := q.db.QueryRow(ctx, createUserCredential, arg.UserID, arg.HashedPassword)
	var i UserCredential
	err := row.Scan(&i.UserID, &i.HashedPassword)
	return i, err
}

const getAllUsers = `-- name: GetAllUsers :many
Select id, username, email, email_verified, created_at from users
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.EmailVerified,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
Select id, username, email, email_verified, created_at from users where email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.EmailVerified,
		&i.CreatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
Select id, username, email, email_verified, created_at from users where id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.EmailVerified,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
Select id, username, email, email_verified, created_at from users where username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.EmailVerified,
		&i.CreatedAt,
	)
	return i, err
}

const getUserCredentialByUserId = `-- name: GetUserCredentialByUserId :one
Select user_id, hashed_password from user_credentials u where u.user_id = $1
`

func (q *Queries) GetUserCredentialByUserId(ctx context.Context, userID uuid.UUID) (UserCredential, error) {
	row := q.db.QueryRow(ctx, getUserCredentialByUserId, userID)
	var i UserCredential
	err := row.Scan(&i.UserID, &i.HashedPassword)
	return i, err
}
