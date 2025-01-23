package auth

import "errors"

var (
	ErrRevokedRefreshToken = errors.New("refresh token is already used")
)
