package auth

import (
	"context"
	"errors"
	"time"

	"ecommerce/internal/domain/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Login(ctx context.Context, req LoginRequest) (string, string, error) {
    email := req.Email
    password := req.Password
    
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil || user == nil {
        return "", "", errors.New("invalid credentials")
    }

    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
        return "", "", errors.New("invalid credentials")
    }

    accessToken, err := s.generateToken(user)
    if err != nil {
        return "", "", err
    }

    refreshToken := uuid.New().String()
    if err := s.userRepo.StoreRefreshToken(ctx, user.ID, refreshToken); err != nil {
        return "", "", err
    }

    return accessToken, refreshToken, nil
}


func (s *service) generateToken(user *entity.User) (string, error) {
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}
