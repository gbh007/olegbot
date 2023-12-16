package domain

import "errors"

var (
	QuoteAlreadyExistsError = errors.New("quote already exists")
	PermissionDeniedError   = errors.New("permission denied")
)
