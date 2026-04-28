package postgres

import (
	"context"
	"ecommerce/internal/domain/entity"
	"ecommerce/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) repository.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	query := `
        INSERT INTO orders (
            id, user_id, status, total_amount,
            shipping_address, payment_method,
            created_at, updated_at
        )
        VALUES (
            :id, :user_id, :status, :total_amount,
            :shipping_address, :payment_method,
            :created_at, :updated_at
        )
    `
	_, err := r.db.NamedExecContext(ctx, query, order)
	return err
}

func (r *orderRepository) CreateOrderItem(ctx context.Context, item *entity.OrderItem) error {
	query := `
        INSERT INTO order_items (id, order_id, product_id, quantity, price)
        VALUES (:id, :order_id, :product_id, :quantity, :price)
    `
	_, err := r.db.NamedExecContext(ctx, query, item)
	return err
}

func (r *orderRepository) GetOrder(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*entity.Order, []*entity.OrderItem, error) {
	var order entity.Order

	err := r.db.GetContext(ctx, &order,
		`SELECT * FROM orders WHERE id = $1 AND user_id = $2`,
		orderID, userID,
	)
	if err != nil {
		return nil, nil, err
	}

	items := []*entity.OrderItem{}
	err = r.db.SelectContext(ctx, &items,
		`SELECT * FROM order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, nil, err
	}

	return &order, items, nil
}

func (r *orderRepository) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*entity.Order, error) {
	orders := []*entity.Order{}
	err := r.db.SelectContext(ctx, &orders,
		`SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	return orders, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	return err
}
