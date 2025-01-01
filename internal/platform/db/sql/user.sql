-- name: CreateUser :one
INSERT INTO users (
  id,
  username,
  email,
  email_verified
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: CreateUserCredential :one
INSERT INTO user_credentials (
  user_id,
  hashed_password,
  salt
) VALUES (
  $1, $2, $3
) RETURNING *;


-- name: GetUserById :one
Select * from users where id = $1;

-- name: GetAllUsers :many
Select * from users;

-- name: GetUserByUsername :one
Select * from users where username = $1;

-- name: GetUserByEmail :one
Select * from users where email = $1;

-- name: GetUserCredentialByUserId :one
Select * from user_credentials u where u.user_id = $1;