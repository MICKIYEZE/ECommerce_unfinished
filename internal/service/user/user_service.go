package user

import (
    "context"
    "fmt"
    "strings"

    "ecommerce/internal/domain"
    "ecommerce/internal/repository/postgrers"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type UserService interface {
    GetProfile(ctx context.Context, id uuid.UUID) (*domain.User, error)
    UpdateProfile(ctx context.Context, id uuid.UUID, req UpdateProfileRequest, isAdmin bool) (*domain.User, error)
    ListUser(ctx context.Context, filter UserFilter) (*UserListResponse, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
    userRepo postgres.UserRepository
}

func NewService(userRepo postgres.UserRepository) UserService {
    return &service{userRepo: userRepo}
}

func (s *service) GetProfile(ctx context.Context, id uuid.UUID) (*domain.User, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, ErrUserNotFound
    }
    return user, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
    if err := s.userRepo.Delete(ctx, id); err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    return nil
}

func (s *service) ListUser(ctx context.Context, filter UserFilter) (*UserListResponse, error) {
    if filter.Limit <= 0 {
        filter.Limit = 20
    }

    if filter.Limit > 100 {
        filter.Limit = 100
    }

    if filter.Offset < 0 {
        filter.Offset = 0
    }

    users, err := s.userRepo.List(ctx, filter.Limit, filter.Offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }

    filtered := []*domain.User{}

    for _, u := range users {
        if filter.Role != nil && u.Role != *filter.Role {
            continue
        }

        if filter.Search != "" &&
            !strings.Contains(strings.ToLower(u.FirstName), strings.ToLower(filter.Search)) &&
            !strings.Contains(strings.ToLower(u.Surname), strings.ToLower(filter.Search)) {
            continue
        }

        filtered = append(filtered, u)
    }

    return &UserListResponse{
        Users: filtered,
        Total: len(filtered),
    }, nil
}

func (s *service) UpdateProfile(ctx context.Context, id uuid.UUID, req UpdateProfileRequest, isAdmin bool) (*domain.User, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, ErrUserNotFound
    }

    if req.FirstName != nil {
        user.FirstName = *req.FirstName
    }

    if req.LastName != nil {
        user.Surname = *req.LastName
    }

    if req.Email != nil {
        user.Email = *req.Email
    }

    if req.Password != nil {
        hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
        if err != nil {
            return nil, fmt.Errorf("failed to hash password: %w", err)
        }
        user.PasswordHash = string(hash)
    }

    if req.Role != nil && isAdmin {
        user.Role = *req.Role
    }

    if err := s.userRepo.Update(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to update user: %w", err)
    }

    return user, nil
}
