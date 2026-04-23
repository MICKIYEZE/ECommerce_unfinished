package repository

import (
    "context"

    "github.com/google/uuid"
    "ecommerce/internal/domain"
)

type OrderRepository interface {
    Create(ctx context.Context, order *domain.Order, items []*domain.OrderItem) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, []*domain.OrderItem, error)
    ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}
