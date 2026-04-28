package cart

import (
    "context"
    "errors"

    "ecommerce/internal/domain/entity"
    "ecommerce/internal/domain/repository"

    "github.com/google/uuid"
)

type CartService interface {
    AddItem(userID uuid.UUID, productID uuid.UUID, quantity int) error
    UpdateItem(userID uuid.UUID, productID uuid.UUID, quantity int) error
    RemoveItem(userID uuid.UUID, productID uuid.UUID) error
    ClearCart(userID uuid.UUID) error
    GetCart(userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error)
}

type service struct {
    cartRepo    repository.CartRepository
    productRepo repository.ProductRepository
}

func NewService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
    return &service{
        cartRepo:    cartRepo,
        productRepo: productRepo,
    }
}

// AddItem adds a product to the user's cart
func (s *service) AddItem(userID uuid.UUID, productID uuid.UUID, quantity int) error {
    if quantity <= 0 {
        return errors.New("quantity must be greater than zero")
    }

    product, err := s.productRepo.GetByID(context.Background(), productID)
    if err != nil || product == nil {
        return errors.New("product not found")
    }

    if product.Stock < quantity {
        return errors.New("not enough stock")
    }

    return s.cartRepo.AddItem(context.Background(), userID, productID, quantity)
}

// UpdateItem updates the quantity of a product in the cart
func (s *service) UpdateItem(userID uuid.UUID, productID uuid.UUID, quantity int) error {
    if quantity <= 0 {
        return errors.New("quantity must be greater than zero")
    }

    product, err := s.productRepo.GetByID(context.Background(), productID)
    if err != nil || product == nil {
        return errors.New("product not found")
    }

    if product.Stock < quantity {
        return errors.New("not enough stock")
    }

    return s.cartRepo.UpdateItem(context.Background(), userID, productID, quantity)
}

// RemoveItem removes a product from the cart
func (s *service) RemoveItem(userID uuid.UUID, productID uuid.UUID) error {
    return s.cartRepo.RemoveItem(context.Background(), userID, productID)
}

// ClearCart removes all items from the user's cart
func (s *service) ClearCart(userID uuid.UUID) error {
    return s.cartRepo.ClearCart(context.Background(), userID)
}

// GetCart returns the cart and its items
func (s *service) GetCart(userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error) {
    return s.cartRepo.GetCart(context.Background(), userID)
}
