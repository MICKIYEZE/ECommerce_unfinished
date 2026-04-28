package auth

import "errors"

var (
    ErrEmailRequired       = errors.New("email is required")
    ErrPasswordRequired    = errors.New("password is required")
    ErrWeakPassword        = errors.New("password is too weak")
    ErrUserAlreadyExists   = errors.New("user already exists")
    ErrInvalidCredentials  = errors.New("invalid email or password")
    ErrInvalidToken        = errors.New("invalid token")
    ErrTokenExpired        = errors.New("token has expired")
)
