package errs

import "errors"

var (
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrInvalidCredentials  = errors.New("invalid credential")
	ErrTokenNotFound       = errors.New("refresh token not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotFound        = errors.New("user not found")

	ErrFileTooLarge     = errors.New("file too large")
	ErrInvalidmediaType = errors.New("invalid media type")

	ErrCacheKeyNotFound = errors.New("key not found")
)
