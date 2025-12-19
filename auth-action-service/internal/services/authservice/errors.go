package authservice

import "errors"

var (
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
)


