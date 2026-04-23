package service

import (
    "context"
    "ecommerce/internal/domain/entity"
)

type RegisterRequest struct {
    Email     string
    Password  string
    FirstName string
    LastName  string
}

type LoginRequest struct {
    Email    string
    Password string
}

type AuthResponse struct {
    Token string
}

type AuthService interface {
    Register(ctx context.Context, req RegisterRequest) (*entity.User, error)
    Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}
