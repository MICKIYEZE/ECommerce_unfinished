package repository

import (
    "context"

    "github.com/google/uuid"
    "ecommerce/internal/domain/entity"
)

type OrderRepository interface {
    CreateOrder(ctx context.Context, order *entity.Order) error
    CreateOrderItem(ctx context.Context, item *entity.OrderItem) error

    GetOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*entity.Order, []*entity.OrderItem, error)
    GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*entity.Order, error)

    UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}
