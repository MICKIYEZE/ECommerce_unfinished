package repository

import (
    "context"

    "github.com/google/uuid"
    "ecommerce/internal/domain/entity"
)

type ProductRepository interface {
    Create(ctx context.Context, product *entity.Product) error
    Update(ctx context.Context, product *entity.Product) error
    Delete(ctx context.Context, id uuid.UUID) error

    GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
    List(ctx context.Context, filter entity.ProductFilter) ([]*entity.Product, error)

    DecrementStock(ctx context.Context, productID uuid.UUID, qty int) error
    IncrementStock(ctx context.Context, productID uuid.UUID, qty int) error
}
