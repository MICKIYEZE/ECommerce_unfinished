package cart

import (
    "context"
    "errors"
    "testing"

    "ecommerce/internal/domain/entity"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

//
// FULL MOCKS
//

// ----------------------
// CartRepository Mock
// ----------------------
type MockCartRepo struct{ mock.Mock }

func (m *MockCartRepo) AddItem(ctx context.Context, userID, productID uuid.UUID, qty int) error {
    args := m.Called(ctx, userID, productID, qty)
    return args.Error(0)
}

func (m *MockCartRepo) UpdateItem(ctx context.Context, userID, productID uuid.UUID, qty int) error {
    args := m.Called(ctx, userID, productID, qty)
    return args.Error(0)
}

func (m *MockCartRepo) RemoveItem(ctx context.Context, userID, productID uuid.UUID) error {
    args := m.Called(ctx, userID, productID)
    return args.Error(0)
}

func (m *MockCartRepo) ClearCart(ctx context.Context, userID uuid.UUID) error {
    args := m.Called(ctx, userID)
    return args.Error(0)
}

func (m *MockCartRepo) GetCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error) {
    args := m.Called(ctx, userID)
    return args.Get(0).(*entity.Cart), args.Get(1).([]*entity.CartItem), args.Error(2)
}

// ----------------------
// ProductRepository Mock
// ----------------------
type MockProductRepo struct{ mock.Mock }

func (m *MockProductRepo) Create(ctx context.Context, p *entity.Product) error {
    args := m.Called(ctx, p)
    return args.Error(0)
}

func (m *MockProductRepo) Update(ctx context.Context, p *entity.Product) error {
    args := m.Called(ctx, p)
    return args.Error(0)
}

func (m *MockProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepo) List(ctx context.Context, filter entity.ProductFilter) ([]*entity.Product, error) {
    args := m.Called(ctx, filter)
    return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *MockProductRepo) DecrementStock(ctx context.Context, id uuid.UUID, qty int) error {
    args := m.Called(ctx, id, qty)
    return args.Error(0)
}

func (m *MockProductRepo) IncrementStock(ctx context.Context, id uuid.UUID, qty int) error {
    args := m.Called(ctx, id, qty)
    return args.Error(0)
}

//
// TESTS
//

// ----------------------
// AddItem
// ----------------------

func TestAddItem_QuantityZero(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    err := svc.AddItem(uuid.New(), uuid.New(), 0)
    assert.EqualError(t, err, "quantity must be greater than zero")
}

func TestAddItem_ProductNotFound(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    productRepo.On("GetByID", mock.Anything, mock.Anything).
        Return((*entity.Product)(nil), errors.New("not found"))

    err := svc.AddItem(uuid.New(), uuid.New(), 2)
    assert.EqualError(t, err, "product not found")
}

func TestAddItem_InsufficientStock(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    productRepo.On("GetByID", mock.Anything, mock.Anything).
        Return(&entity.Product{Stock: 1}, nil)

    err := svc.AddItem(uuid.New(), uuid.New(), 5)
    assert.EqualError(t, err, "not enough stock")
}

func TestAddItem_RepoError(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    productRepo.On("GetByID", mock.Anything, mock.Anything).
        Return(&entity.Product{Stock: 10}, nil)

    cartRepo.On("AddItem", mock.Anything, mock.Anything, mock.Anything, 3).
        Return(errors.New("db error"))

    err := svc.AddItem(uuid.New(), uuid.New(), 3)
    assert.EqualError(t, err, "db error")
}

func TestAddItem_Success(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    productRepo.On("GetByID", mock.Anything, mock.Anything).
        Return(&entity.Product{Stock: 10}, nil)

    cartRepo.On("AddItem", mock.Anything, mock.Anything, mock.Anything, 2).
        Return(nil)

    err := svc.AddItem(uuid.New(), uuid.New(), 2)
    assert.NoError(t, err)
}

// ----------------------
// UpdateItem
// ----------------------

func TestUpdateItem_QuantityZero(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    err := svc.UpdateItem(uuid.New(), uuid.New(), 0)
    assert.EqualError(t, err, "quantity must be greater than zero")
}

func TestUpdateItem_ProductNotFound(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    productRepo.On("GetByID", mock.Anything, mock.Anything).
        Return((*entity.Product)(nil), errors.New("not found"))

    err := svc.UpdateItem(uuid.New(), uuid.New(), 2)
    assert.EqualError(t, err, "product not found")
}

func TestUpdateItem_InsufficientStock(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    productRepo.On("GetByID", mock.Anything, mock.Anything).
        Return(&entity.Product{Stock: 1}, nil)

    err := svc.UpdateItem(uuid.New(), uuid.New(), 5)
    assert.EqualError(t, err, "not enough stock")
}

func TestUpdateItem_RepoError(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    productRepo.On("GetByID", mock.Anything, mock.Anything).
        Return(&entity.Product{Stock: 10}, nil)

    cartRepo.On("UpdateItem", mock.Anything, mock.Anything, mock.Anything, 3).
        Return(errors.New("db error"))

    err := svc.UpdateItem(uuid.New(), uuid.New(), 3)
    assert.EqualError(t, err, "db error")
}

func TestUpdateItem_Success(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    productRepo.On("GetByID", mock.Anything, mock.Anything).
        Return(&entity.Product{Stock: 10}, nil)

    cartRepo.On("UpdateItem", mock.Anything, mock.Anything, mock.Anything, 2).
        Return(nil)

    err := svc.UpdateItem(uuid.New(), uuid.New(), 2)
    assert.NoError(t, err)
}

// ----------------------
// RemoveItem
// ----------------------

func TestRemoveItem_Error(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    cartRepo.On("RemoveItem", mock.Anything, mock.Anything, mock.Anything).
        Return(errors.New("db error"))

    err := svc.RemoveItem(uuid.New(), uuid.New())
    assert.EqualError(t, err, "db error")
}

func TestRemoveItem_Success(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    cartRepo.On("RemoveItem", mock.Anything, mock.Anything, mock.Anything).
        Return(nil)

    err := svc.RemoveItem(uuid.New(), uuid.New())
    assert.NoError(t, err)
}

// ----------------------
// ClearCart
// ----------------------

func TestClearCart_Error(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    cartRepo.On("ClearCart", mock.Anything, mock.Anything).
        Return(errors.New("db error"))

    err := svc.ClearCart(uuid.New())
    assert.EqualError(t, err, "db error")
}

func TestClearCart_Success(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    cartRepo.On("ClearCart", mock.Anything, mock.Anything).
        Return(nil)

    err := svc.ClearCart(uuid.New())
    assert.NoError(t, err)
}

// ----------------------
// GetCart
// ----------------------

func TestGetCart_Error(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    cartRepo.On("GetCart", mock.Anything, mock.Anything).
        Return(&entity.Cart{}, []*entity.CartItem{}, errors.New("db error"))

    _, _, err := svc.GetCart(uuid.New())
    assert.EqualError(t, err, "db error")
}

func TestGetCart_Success(t *testing.T) {
    cartRepo := new(MockCartRepo)
    productRepo := new(MockProductRepo)
    svc := NewService(cartRepo, productRepo)

    cart := &entity.Cart{}
    items := []*entity.CartItem{}

    cartRepo.On("GetCart", mock.Anything, mock.Anything).
        Return(cart, items, nil)

    c, i, err := svc.GetCart(uuid.New())

    assert.NoError(t, err)
    assert.Equal(t, cart, c)
    assert.Equal(t, items, i)
}
