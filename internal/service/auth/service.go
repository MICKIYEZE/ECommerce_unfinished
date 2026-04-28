package auth

import (
    "context"

    "ecommerce/internal/domain/entity"
    "ecommerce/internal/repository/postgres"
)

type AuthService interface {
    Register(ctx context.Context, req RegisterRequest) (*entity.User, error)
    Login(ctx context.Context, req LoginRequest) (string, string, error)
    ValidateToken(tokenString string) (*Claims, error)
    RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

type service struct {
    userRepo  postgres.UserRepository
    jwtSecret []byte
}

func NewService(userRepo postgres.UserRepository, jwtSecret string) AuthService {
    return &service{
        userRepo:  userRepo,
        jwtSecret: []byte(jwtSecret),
    }
}
