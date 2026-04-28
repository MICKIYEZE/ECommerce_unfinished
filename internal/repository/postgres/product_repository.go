package postgres

import (
    "context"
    "ecommerce/internal/domain/entity"
    "ecommerce/internal/domain/repository"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
)

type productRepository struct {
    db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) repository.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *entity.Product) error {
    query := `
        INSERT INTO products (id, name, description, price, stock, created_at, updated_at)
        VALUES (:id, :name, :description, :price, :stock, :created_at, :updated_at)
    `
    _, err := r.db.NamedExecContext(ctx, query, p)
    return err
}

func (r *productRepository) Update(ctx context.Context, p *entity.Product) error {
    query := `
        UPDATE products
        SET name = :name,
            description = :description,
            price = :price,
            stock = :stock,
            updated_at = :updated_at
        WHERE id = :id
    `
    _, err := r.db.NamedExecContext(ctx, query, p)
    return err
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
    _, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
    return err
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
    var p entity.Product
    err := r.db.GetContext(ctx, &p, `SELECT * FROM products WHERE id = $1`, id)
    if err != nil {
        return nil, err
    }
    return &p, nil
}

func (r *productRepository) List(ctx context.Context, filter entity.ProductFilter) ([]*entity.Product, error) {
    products := []*entity.Product{}
    query := `SELECT * FROM products ORDER BY created_at DESC`
    err := r.db.SelectContext(ctx, &products, query)
    return products, err
}

func (r *productRepository) DecrementStock(ctx context.Context, productID uuid.UUID, qty int) error {
    _, err := r.db.ExecContext(ctx,
        `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`,
        qty, productID,
    )
    return err
}

func (r *productRepository) IncrementStock(ctx context.Context, productID uuid.UUID, qty int) error {
    _, err := r.db.ExecContext(ctx,
        `UPDATE products SET stock = stock + $1 WHERE id = $2`,
        qty, productID,
    )
    return err
}
