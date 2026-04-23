package repository

import (
    "context"

    "github.com/google/uuid"
    "ecommerce/internal/domain"
)

type ProductRepository interface {
    Create(ctx context.Context, product *domain.Product) error
    Update(ctx context.Context, product *domain.Product) error
    Delete(ctx context.Context, id uuid.UUID) error

    GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
    List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, error)

    DecreaseStock(ctx context.Context, productID uuid.UUID, qty int) error
    IncreaseStock(ctx context.Context, productID uuid.UUID, qty int) error
}
