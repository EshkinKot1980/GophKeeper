package errors

import "errors"

var (
	ErrUnexpected             = errors.New("unexpected error")
	ErrAuthUserAlreadyExists  = errors.New("user already exists")
	ErrAuthInvalidCredentials = errors.New("invalid credentials")
	ErrAuthInvalidToken       = errors.New("invalid token")
	ErrAuthTokenExpired       = errors.New("token expired")
	ErrSecretInalidData       = errors.New("invalid secret data")
)
