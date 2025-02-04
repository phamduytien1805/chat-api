package common

import "errors"

var (
	ErrUserNotFound           = errors.New("user_not_found")
	ErrorUserResourceConflict = errors.New("user_resource_conflict")
)
