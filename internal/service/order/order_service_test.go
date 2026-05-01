package order_test

import (
    "context"
    "errors"
    "testing"

    "ecommerce/internal/domain/entity"
    "ecommerce/internal/service/order"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

//
// FULL MOCKS — implement ENTIRE interfaces
//

// ----------------------
// Mock OrderRepository
// ----------------------
type MockOrderRepo struct {
    CreateOrderFn     func(ctx context.Context, o *entity.Order) error
    CreateOrderItemFn func(ctx context.Context, item *entity.OrderItem) error
    GetOrderFn        func(ctx context.Context, userID, orderID uuid.UUID) (*entity.Order, []*entity.OrderItem, error)
    GetUserOrdersFn   func(ctx context.Context, userID uuid.UUID) ([]*entity.Order, error)
    UpdateStatusFn    func(ctx context.Context, id uuid.UUID, status string) error
}

func (m *MockOrderRepo) CreateOrder(ctx context.Context, o *entity.Order) error {
    return m.CreateOrderFn(ctx, o)
}
func (m *MockOrderRepo) CreateOrderItem(ctx context.Context, item *entity.OrderItem) error {
    return m.CreateOrderItemFn(ctx, item)
}
func (m *MockOrderRepo) GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*entity.Order, []*entity.OrderItem, error) {
    return m.GetOrderFn(ctx, userID, orderID)
}
func (m *MockOrderRepo) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*entity.Order, error) {
    return m.GetUserOrdersFn(ctx, userID)
}
func (m *MockOrderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
    return m.UpdateStatusFn(ctx, id, status)
}

// ----------------------
// Mock CartRepository
// ----------------------
type MockCartRepo struct {
    GetCartFn   func(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error)
    ClearCartFn func(ctx context.Context, userID uuid.UUID) error

    // unused but required by interface
    AddItemFn    func(ctx context.Context, userID, productID uuid.UUID, qty int) error
    UpdateItemFn func(ctx context.Context, userID, productID uuid.UUID, qty int) error
    RemoveItemFn func(ctx context.Context, userID, productID uuid.UUID) error
}

func (m *MockCartRepo) GetCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error) {
    return m.GetCartFn(ctx, userID)
}
func (m *MockCartRepo) ClearCart(ctx context.Context, userID uuid.UUID) error {
    return m.ClearCartFn(ctx, userID)
}
func (m *MockCartRepo) AddItem(ctx context.Context, userID, productID uuid.UUID, qty int) error {
    if m.AddItemFn != nil {
        return m.AddItemFn(ctx, userID, productID, qty)
    }
    return nil
}
func (m *MockCartRepo) UpdateItem(ctx context.Context, userID, productID uuid.UUID, qty int) error {
    if m.UpdateItemFn != nil {
        return m.UpdateItemFn(ctx, userID, productID, qty)
    }
    return nil
}
func (m *MockCartRepo) RemoveItem(ctx context.Context, userID, productID uuid.UUID) error {
    if m.RemoveItemFn != nil {
        return m.RemoveItemFn(ctx, userID, productID)
    }
    return nil
}

// ----------------------
// Mock ProductRepository
// ----------------------
type MockProductRepo struct {
    GetByIDFn        func(ctx context.Context, id uuid.UUID) (*entity.Product, error)
    DecrementStockFn func(ctx context.Context, id uuid.UUID, qty int) error
    IncrementStockFn func(ctx context.Context, id uuid.UUID, qty int) error

    // unused but required
    CreateFn func(ctx context.Context, p *entity.Product) error
    UpdateFn func(ctx context.Context, p *entity.Product) error
    DeleteFn func(ctx context.Context, id uuid.UUID) error
    ListFn   func(ctx context.Context, filter entity.ProductFilter) ([]*entity.Product, error)
}

func (m *MockProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
    return m.GetByIDFn(ctx, id)
}
func (m *MockProductRepo) DecrementStock(ctx context.Context, id uuid.UUID, qty int) error {
    return m.DecrementStockFn(ctx, id, qty)
}
func (m *MockProductRepo) IncrementStock(ctx context.Context, id uuid.UUID, qty int) error {
    return m.IncrementStockFn(ctx, id, qty)
}
func (m *MockProductRepo) Create(ctx context.Context, p *entity.Product) error {
    if m.CreateFn != nil {
        return m.CreateFn(ctx, p)
    }
    return nil
}
func (m *MockProductRepo) Update(ctx context.Context, p *entity.Product) error {
    if m.UpdateFn != nil {
        return m.UpdateFn(ctx, p)
    }
    return nil
}
func (m *MockProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
    if m.DeleteFn != nil {
        return m.DeleteFn(ctx, id)
    }
    return nil
}
func (m *MockProductRepo) List(ctx context.Context, filter entity.ProductFilter) ([]*entity.Product, error) {
    if m.ListFn != nil {
        return m.ListFn(ctx, filter)
    }
    return nil, nil
}

//
// TESTS
//

func TestCreateOrder_Success(t *testing.T) {
    userID := uuid.New()
    productID := uuid.New()

    mockOrderRepo := &MockOrderRepo{
        CreateOrderFn:     func(ctx context.Context, o *entity.Order) error { return nil },
        CreateOrderItemFn: func(ctx context.Context, item *entity.OrderItem) error { return nil },
    }
    mockCartRepo := &MockCartRepo{
        GetCartFn: func(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error) {
            return &entity.Cart{UserID: userID}, []*entity.CartItem{
                {ProductID: productID, Quantity: 2},
            }, nil
        },
        ClearCartFn: func(ctx context.Context, userID uuid.UUID) error { return nil },
    }
    mockProductRepo := &MockProductRepo{
        GetByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
            return &entity.Product{ID: id, Name: "Test Product", Price: 10, Stock: 10}, nil
        },
        DecrementStockFn: func(ctx context.Context, id uuid.UUID, qty int) error { return nil },
    }

    svc := order.NewService(mockOrderRepo, mockCartRepo, mockProductRepo)

    addr := entity.Address{Street: "Test", City: "City"}

    o, items, err := svc.CreateOrder(userID, addr, "card")

    assert.NoError(t, err)
    assert.NotNil(t, o)
    assert.Len(t, items, 1)
    assert.Equal(t, float64(20), o.TotalAmount)
}

func TestCreateOrder_EmptyCart(t *testing.T) {
    mockOrderRepo := &MockOrderRepo{}
    mockCartRepo := &MockCartRepo{
        GetCartFn: func(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error) {
            return &entity.Cart{}, []*entity.CartItem{}, nil
        },
    }
    mockProductRepo := &MockProductRepo{}

    svc := order.NewService(mockOrderRepo, mockCartRepo, mockProductRepo)

    _, _, err := svc.CreateOrder(uuid.New(), entity.Address{}, "card")

    assert.EqualError(t, err, "cart is empty")
}

func TestCreateOrder_ProductNotFound(t *testing.T) {
    productID := uuid.New()

    mockOrderRepo := &MockOrderRepo{}
    mockCartRepo := &MockCartRepo{
        GetCartFn: func(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error) {
            return &entity.Cart{}, []*entity.CartItem{
                {ProductID: productID, Quantity: 1},
            }, nil
        },
    }
    mockProductRepo := &MockProductRepo{
        GetByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
            return nil, errors.New("not found")
        },
    }

    svc := order.NewService(mockOrderRepo, mockCartRepo, mockProductRepo)

    _, _, err := svc.CreateOrder(uuid.New(), entity.Address{}, "card")

    assert.EqualError(t, err, "product not found")
}

func TestCreateOrder_NotEnoughStock(t *testing.T) {
    productID := uuid.New()

    mockOrderRepo := &MockOrderRepo{}
    mockCartRepo := &MockCartRepo{
        GetCartFn: func(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error) {
            return &entity.Cart{}, []*entity.CartItem{
                {ProductID: productID, Quantity: 5},
            }, nil
        },
    }
    mockProductRepo := &MockProductRepo{
        GetByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
            return &entity.Product{ID: id, Name: "Test", Price: 10, Stock: 1}, nil
        },
    }

    svc := order.NewService(mockOrderRepo, mockCartRepo, mockProductRepo)

    _, _, err := svc.CreateOrder(uuid.New(), entity.Address{}, "card")

    assert.EqualError(t, err, "not enough stock for product: Test")
}

func TestCancelOrder_Success(t *testing.T) {
    orderID := uuid.New()
    userID := uuid.New()

    mockOrderRepo := &MockOrderRepo{
        GetOrderFn: func(ctx context.Context, uid, oid uuid.UUID) (*entity.Order, []*entity.OrderItem, error) {
            return &entity.Order{ID: oid, Status: "pending"}, []*entity.OrderItem{
                {ProductID: uuid.New(), Quantity: 2},
            }, nil
        },
        UpdateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error { return nil },
    }
    mockCartRepo := &MockCartRepo{}
    mockProductRepo := &MockProductRepo{
        IncrementStockFn: func(ctx context.Context, id uuid.UUID, qty int) error { return nil },
    }

    svc := order.NewService(mockOrderRepo, mockCartRepo, mockProductRepo)

    err := svc.CancelOrder(userID, orderID)

    assert.NoError(t, err)
}

func TestCancelOrder_NotPending(t *testing.T) {
    orderID := uuid.New()
    userID := uuid.New()

    mockOrderRepo := &MockOrderRepo{
        GetOrderFn: func(ctx context.Context, uid, oid uuid.UUID) (*entity.Order, []*entity.OrderItem, error) {
            return &entity.Order{ID: oid, Status: "completed"}, nil, nil
        },
    }
    mockCartRepo := &MockCartRepo{}
    mockProductRepo := &MockProductRepo{}

    svc := order.NewService(mockOrderRepo, mockCartRepo, mockProductRepo)

    err := svc.CancelOrder(userID, orderID)

    assert.EqualError(t, err, "only pending orders can be cancelled")
}

func TestGetOrder(t *testing.T) {
    orderID := uuid.New()
    userID := uuid.New()

    mockOrderRepo := &MockOrderRepo{
        GetOrderFn: func(ctx context.Context, uid, oid uuid.UUID) (*entity.Order, []*entity.OrderItem, error) {
            return &entity.Order{ID: oid}, []*entity.OrderItem{}, nil
        },
    }
    mockCartRepo := &MockCartRepo{}
    mockProductRepo := &MockProductRepo{}

    svc := order.NewService(mockOrderRepo, mockCartRepo, mockProductRepo)

    o, items, err := svc.GetOrder(userID, orderID)

    assert.NoError(t, err)
    assert.NotNil(t, o)
    assert.NotNil(t, items)
}

func TestGetUserOrders(t *testing.T) {
    userID := uuid.New()

    mockOrderRepo := &MockOrderRepo{
        GetUserOrdersFn: func(ctx context.Context, uid uuid.UUID) ([]*entity.Order, error) {
            return []*entity.Order{{ID: uuid.New()}}, nil
        },
    }
    mockCartRepo := &MockCartRepo{}
    mockProductRepo := &MockProductRepo{}

    svc := order.NewService(mockOrderRepo, mockCartRepo, mockProductRepo)

    orders, err := svc.GetUserOrders(userID)

    assert.NoError(t, err)
    assert.Len(t, orders, 1)
}
