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
  hashed_password
) VALUES (
  $1, $2
) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
  username = COALESCE($2, username),
  email = COALESCE($3, email),
  email_verified = COALESCE($4, email_verified)
WHERE id = $1 RETURNING *;

-- name: UpdateUserByEmail :one
UPDATE users SET
  username = COALESCE($2, username),
  email_verified = COALESCE($3, email_verified)
WHERE email = $1 RETURNING *;


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

-- name: GetUserByEmailOrUsername :one
Select * from users where email = $1 OR username = $1;
