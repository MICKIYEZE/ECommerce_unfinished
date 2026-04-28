package order

import (
    "context"
    "errors"
    "time"

    "ecommerce/internal/domain/entity"
    "ecommerce/internal/domain/repository"

    "github.com/google/uuid"
)

type OrderService interface {
    CreateOrder(userID uuid.UUID, address entity.Address, paymentMethod string) (*entity.Order, []*entity.OrderItem, error)
    GetOrder(userID uuid.UUID, orderID uuid.UUID) (*entity.Order, []*entity.OrderItem, error)
    GetUserOrders(userID uuid.UUID) ([]*entity.Order, error)
    CancelOrder(userID uuid.UUID, orderID uuid.UUID) error
}

type service struct {
    orderRepo   repository.OrderRepository
    cartRepo    repository.CartRepository
    productRepo repository.ProductRepository
}

func NewService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository) OrderService {
    return &service{
        orderRepo:   orderRepo,
        cartRepo:    cartRepo,
        productRepo: productRepo,
    }
}

func (s *service) CreateOrder(userID uuid.UUID, address entity.Address, paymentMethod string) (*entity.Order, []*entity.OrderItem, error) {
    ctx := context.Background()

    cart, items, err := s.cartRepo.GetCart(ctx, userID)
    if err != nil {
        return nil, nil, err
    }

    if cart == nil || len(items) == 0 {
        return nil, nil, errors.New("cart is empty")
    }

    var total float64

    for _, item := range items {
        product, err := s.productRepo.GetByID(ctx, item.ProductID)
        if err != nil || product == nil {
            return nil, nil, errors.New("product not found")
        }
        if product.Stock < item.Quantity {
            return nil, nil, errors.New("not enough stock for product: " + product.Name)
        }

        total += product.Price * float64(item.Quantity)
    }

    // Create order
    order := &entity.Order{
        ID:              uuid.New(),
        UserID:          userID,
        Status:          "pending",
        TotalAmount:     total,
        ShippingAddress: address,
        PaymentMethod:   paymentMethod,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }

    if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
        return nil, nil, err
    }

    // Create order items + decrement stock
    var orderItems []*entity.OrderItem

    for _, item := range items {
        product, _ := s.productRepo.GetByID(ctx, item.ProductID)

        orderItem := &entity.OrderItem{
            ID:        uuid.New(),
            OrderID:   order.ID,
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            Price:     product.Price,
        }

        if err := s.orderRepo.CreateOrderItem(ctx, orderItem); err != nil {
            return nil, nil, err
        }

        orderItems = append(orderItems, orderItem)

        // Decrement stock
        if err := s.productRepo.DecrementStock(ctx, item.ProductID, item.Quantity); err != nil {
            return nil, nil, err
        }
    }

    // Clear cart
    _ = s.cartRepo.ClearCart(ctx, userID)

    return order, orderItems, nil
}

// GetOrder returns an order and its items
func (s *service) GetOrder(userID uuid.UUID, orderID uuid.UUID) (*entity.Order, []*entity.OrderItem, error) {
    return s.orderRepo.GetOrder(context.Background(), userID, orderID)
}

// GetUserOrders returns all orders for a user
func (s *service) GetUserOrders(userID uuid.UUID) ([]*entity.Order, error) {
    return s.orderRepo.GetUserOrders(context.Background(), userID)
}

// CancelOrder cancels an order and restores stock
func (s *service) CancelOrder(userID uuid.UUID, orderID uuid.UUID) error {
    ctx := context.Background()

    order, items, err := s.orderRepo.GetOrder(ctx, userID, orderID)
    if err != nil {
        return err
    }

    if order.Status != "pending" {
        return errors.New("only pending orders can be cancelled")
    }

    // Restore stock
    for _, item := range items {
        if err := s.productRepo.IncrementStock(ctx, item.ProductID, item.Quantity); err != nil {
            return err
        }
    }

    // Update order status
    return s.orderRepo.UpdateStatus(ctx, orderID, "cancelled")
}
