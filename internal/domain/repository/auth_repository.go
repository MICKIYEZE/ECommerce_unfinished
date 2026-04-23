package repository

import (
    "context"

    "github.com/google/uuid"
    "ecommerce/internal/domain/entity"
)

type AuthRepository interface {
    StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string) error
    GetUserByRefreshToken(ctx context.Context, token string) (*entity.User, error)
    DeleteRefreshToken(ctx context.Context, token string) error
}
