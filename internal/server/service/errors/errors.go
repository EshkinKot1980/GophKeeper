package errors

import "errors"

var (
	ErrUnexpected             = errors.New("unexpected error")
	ErrAuthUserAlreadyExists  = errors.New("user already exists")
	ErrAuthInvalidCredentials = errors.New("invalid credentials")
	ErrAuthInvalidToken       = errors.New("invalid token")
	ErrAuthTokenExpired       = errors.New("token expired")
	ErrSecretInvalidData      = errors.New("invalid secret data")
	ErrSecretNotFound         = errors.New("secret not found")
)
