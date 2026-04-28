package postgres

import (
    "context"
    "ecommerce/internal/domain/entity"
    "ecommerce/internal/domain/repository"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
)

type cartRepository struct {
    db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) repository.CartRepository {
    return &cartRepository{db: db}
}

func (r *cartRepository) AddItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, qty int) error {
    query := `
        INSERT INTO cart_items (user_id, product_id, quantity)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, product_id)
        DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
    `
    _, err := r.db.ExecContext(ctx, query, userID, productID, qty)
    return err
}

func (r *cartRepository) UpdateItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, qty int) error {
    _, err := r.db.ExecContext(ctx,
        `UPDATE cart_items SET quantity = $1 WHERE user_id = $2 AND product_id = $3`,
        qty, userID, productID,
    )
    return err
}

func (r *cartRepository) RemoveItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
    _, err := r.db.ExecContext(ctx,
        `DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`,
        userID, productID,
    )
    return err
}

func (r *cartRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
    _, err := r.db.ExecContext(ctx,
        `DELETE FROM cart_items WHERE user_id = $1`,
        userID,
    )
    return err
}

func (r *cartRepository) GetCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, error) {
    var cart entity.Cart
    cart.UserID = userID

    items := []*entity.CartItem{}
    err := r.db.SelectContext(ctx, &items,
        `SELECT * FROM cart_items WHERE user_id = $1`,
        userID,
    )
    if err != nil {
        return nil, nil, err
    }

    return &cart, items, nil
}
