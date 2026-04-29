package product

import (
	"context"
	"errors"

	"ecommerce/internal/domain/entity"
	"ecommerce/internal/domain/repository"

	"github.com/google/uuid"
)

type ProductService interface {
    Create(product *entity.Product) error
    Update(product *entity.Product) error
    Delete(id uuid.UUID) error
    GetByID(id uuid.UUID) (*entity.Product, error)
    List(filter entity.ProductFilter) ([]*entity.Product, error)
    DecrementStock(productID uuid.UUID, qty int) error
    IncrementStock(productID uuid.UUID, qty int) error
}

type service struct {
    productRepo repository.ProductRepository
}

func NewService(productRepo repository.ProductRepository) ProductService {
    return &service{
        productRepo: productRepo,
    }
}

// Create a new product
func (s *service) Create(product *entity.Product) error {
    if product.Name == "" {
        return errors.New("product name is required")
    }
    if product.Price <= 0 {
        return errors.New("product price must be greater than zero")
    }
    return s.productRepo.Create(context.Background(), product)
}

// Update an existing product
func (s *service) Update(product *entity.Product) error {
    if product.ID == uuid.Nil {
        return errors.New("product ID is required")
    }
    return s.productRepo.Update(context.Background(), product)
}

// Delete a product
func (s *service) Delete(id uuid.UUID) error {
    return s.productRepo.Delete(context.Background(), id)
}

// GetByID returns a single product
func (s *service) GetByID(id uuid.UUID) (*entity.Product, error) {
    return s.productRepo.GetByID(context.Background(), id)
}

// List returns products based on filter
func (s *service) List(filter entity.ProductFilter) ([]*entity.Product, error) {
    return s.productRepo.List(context.Background(), filter)
}

// DecrementStock decreases product stock
func (s *service) DecrementStock(productID uuid.UUID, qty int) error {
    if qty <= 0 {
        return errors.New("quantity must be greater than zero")
    }
    return s.productRepo.DecrementStock(context.Background(), productID, qty)
}

// IncrementStock increases product stock
func (s *service) IncrementStock(productID uuid.UUID, qty int) error {
    if qty <= 0 {
        return errors.New("quantity must be greater than zero")
    }
    return s.productRepo.IncrementStock(context.Background(), productID, qty)
}
