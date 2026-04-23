package service

import (
    "context"
    "ecommerce/internal/domain/entity"
)

type CartService interface {
    AddItem(ctx context.Context, userID string, productID string, qty int) error
    RemoveItem(ctx context.Context, userID string, productID string) error
    GetCart(ctx context.Context, userID string) (*entity.Cart, error)
}
