package product_test

import (
    "context"
    "testing"

    "ecommerce/internal/domain/entity"
    "ecommerce/internal/service/product"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)



type MockProductRepo struct {
    CreateFn         func(ctx context.Context, p *entity.Product) error
    UpdateFn         func(ctx context.Context, p *entity.Product) error
    DeleteFn         func(ctx context.Context, id uuid.UUID) error
    GetByIDFn        func(ctx context.Context, id uuid.UUID) (*entity.Product, error)
    ListFn           func(ctx context.Context, filter entity.ProductFilter) ([]*entity.Product, error)
    DecrementStockFn func(ctx context.Context, id uuid.UUID, qty int) error
    IncrementStockFn func(ctx context.Context, id uuid.UUID, qty int) error
}

func (m *MockProductRepo) Create(ctx context.Context, p *entity.Product) error {
    return m.CreateFn(ctx, p)
}
func (m *MockProductRepo) Update(ctx context.Context, p *entity.Product) error {
    return m.UpdateFn(ctx, p)
}
func (m *MockProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
    return m.DeleteFn(ctx, id)
}
func (m *MockProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
    return m.GetByIDFn(ctx, id)
}
func (m *MockProductRepo) List(ctx context.Context, filter entity.ProductFilter) ([]*entity.Product, error) {
    return m.ListFn(ctx, filter)
}
func (m *MockProductRepo) DecrementStock(ctx context.Context, id uuid.UUID, qty int) error {
    return m.DecrementStockFn(ctx, id, qty)
}
func (m *MockProductRepo) IncrementStock(ctx context.Context, id uuid.UUID, qty int) error {
    return m.IncrementStockFn(ctx, id, qty)
}

//
// TESTS
//

func TestCreateProduct_Success(t *testing.T) {
    mockRepo := &MockProductRepo{
        CreateFn: func(ctx context.Context, p *entity.Product) error { return nil },
    }

    svc := product.NewService(mockRepo)

    p := &entity.Product{
        Name:  "Laptop",
        Price: 1000,
    }

    err := svc.Create(p)

    assert.NoError(t, err)
}

func TestCreateProduct_MissingName(t *testing.T) {
    mockRepo := &MockProductRepo{}
    svc := product.NewService(mockRepo)

    p := &entity.Product{
        Name:  "",
        Price: 100,
    }

    err := svc.Create(p)

    assert.EqualError(t, err, "product name is required")
}

func TestCreateProduct_InvalidPrice(t *testing.T) {
    mockRepo := &MockProductRepo{}
    svc := product.NewService(mockRepo)

    p := &entity.Product{
        Name:  "Laptop",
        Price: 0,
    }

    err := svc.Create(p)

    assert.EqualError(t, err, "product price must be greater than zero")
}

func TestUpdateProduct_Success(t *testing.T) {
    mockRepo := &MockProductRepo{
        UpdateFn: func(ctx context.Context, p *entity.Product) error { return nil },
    }

    svc := product.NewService(mockRepo)

    p := &entity.Product{
        ID:    uuid.New(),
        Name:  "Updated",
        Price: 200,
    }

    err := svc.Update(p)

    assert.NoError(t, err)
}

func TestUpdateProduct_MissingID(t *testing.T) {
    mockRepo := &MockProductRepo{}
    svc := product.NewService(mockRepo)

    p := &entity.Product{
        ID:    uuid.Nil,
        Name:  "Updated",
        Price: 200,
    }

    err := svc.Update(p)

    assert.EqualError(t, err, "product ID is required")
}

func TestDeleteProduct(t *testing.T) {
    mockRepo := &MockProductRepo{
        DeleteFn: func(ctx context.Context, id uuid.UUID) error { return nil },
    }

    svc := product.NewService(mockRepo)

    err := svc.Delete(uuid.New())

    assert.NoError(t, err)
}

func TestGetByID(t *testing.T) {
    id := uuid.New()

    mockRepo := &MockProductRepo{
        GetByIDFn: func(ctx context.Context, pid uuid.UUID) (*entity.Product, error) {
            return &entity.Product{ID: pid, Name: "Test"}, nil
        },
    }

    svc := product.NewService(mockRepo)

    p, err := svc.GetByID(id)

    assert.NoError(t, err)
    assert.Equal(t, id, p.ID)
}

func TestListProducts(t *testing.T) {
    mockRepo := &MockProductRepo{
        ListFn: func(ctx context.Context, filter entity.ProductFilter) ([]*entity.Product, error) {
            return []*entity.Product{{Name: "A"}, {Name: "B"}}, nil
        },
    }

    svc := product.NewService(mockRepo)

    products, err := svc.List(entity.ProductFilter{})

    assert.NoError(t, err)
    assert.Len(t, products, 2)
}

func TestDecrementStock_Success(t *testing.T) {
    mockRepo := &MockProductRepo{
        DecrementStockFn: func(ctx context.Context, id uuid.UUID, qty int) error { return nil },
    }

    svc := product.NewService(mockRepo)

    err := svc.DecrementStock(uuid.New(), 5)

    assert.NoError(t, err)
}

func TestDecrementStock_InvalidQty(t *testing.T) {
    mockRepo := &MockProductRepo{}
    svc := product.NewService(mockRepo)

    err := svc.DecrementStock(uuid.New(), 0)

    assert.EqualError(t, err, "quantity must be greater than zero")
}

func TestIncrementStock_Success(t *testing.T) {
    mockRepo := &MockProductRepo{
        IncrementStockFn: func(ctx context.Context, id uuid.UUID, qty int) error { return nil },
    }

    svc := product.NewService(mockRepo)

    err := svc.IncrementStock(uuid.New(), 3)

    assert.NoError(t, err)
}

func TestIncrementStock_InvalidQty(t *testing.T) {
    mockRepo := &MockProductRepo{}
    svc := product.NewService(mockRepo)

    err := svc.IncrementStock(uuid.New(), 0)

    assert.EqualError(t, err, "quantity must be greater than zero")
}
