package service

import (
    "context"
    "ecommerce/internal/domain/entity"
)

type OrderService interface {
    CreateOrder(ctx context.Context, userID string) (*entity.Order, error)
    GetOrder(ctx context.Context, orderID string) (*entity.Order, error)
    ListOrders(ctx context.Context, userID string) ([]*entity.Order, error)
}
