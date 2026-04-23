package service

import (
    "context"
    "ecommerce/internal/domain/entity"
)

type CreateProductRequest struct {
    Name        string
    Description string
    Price       float64
}

type ProductService interface {
    Create(ctx context.Context, req CreateProductRequest) (*entity.Product, error)
    Get(ctx context.Context, id string) (*entity.Product, error)
    List(ctx context.Context) ([]*entity.Product, error)
}
