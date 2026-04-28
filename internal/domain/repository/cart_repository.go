package repository

import (
    "context"

    "github.com/google/uuid"
    "ecommerce/internal/domain/entity"
)

type CartRepository interface {
    AddItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, qty int) error
    UpdateItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, qty int) error
    RemoveItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
    ClearCart(ctx context.Context, userID uuid.UUID) error
    GetCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error)
}
