package auth

import (
	"context"
	"errors"
	"time"

	"ecommerce/internal/domain/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Register(ctx context.Context, req RegisterRequest) (*entity.User, error) {
    email := req.Email
    password := req.Password
    role := "user"
    
    existing, _ := s.userRepo.GetByEmail(context.Background(), email)
    if existing != nil {
        return nil, errors.New("email already registered")
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &entity.User{
        ID:        uuid.New(),
        Email:     email,
        PasswordHash:  string(hashed),
        Role:      role,
        CreatedAt: time.Now(),
    }

    if err := s.userRepo.Create(context.Background(), user); err != nil {
        return nil, err
    }

    return user, nil
}
