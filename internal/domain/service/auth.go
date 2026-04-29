package service

import (
    "context"
    "ecommerce/internal/domain/entity"
)

type RegisterRequest struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
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
