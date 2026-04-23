package repository

import (
    "context"

    "github.com/google/uuid"
    "ecommerce/internal/domain/entity"
)

type CategoryRepository interface {
    Create(ctx context.Context, category *entity.Category) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
    GetByName(ctx context.Context, name string) (*entity.Category, error)
    List(ctx context.Context) ([]*entity.Category, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
