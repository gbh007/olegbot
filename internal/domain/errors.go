package domain

import "errors"

var (
	AlreadyExistsError    = errors.New("already exists")
	PermissionDeniedError = errors.New("permission denied")
)
