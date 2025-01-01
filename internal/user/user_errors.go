package user

import "errors"

var (
	ErrorUserResourceConflict    = errors.New("username or email are used")
	ErrorUserInvalidAuthenticate = errors.New("username/email or password are incorrect")
)
