package entity

import (
    "time"

    "github.com/google/uuid"
)

type Order struct {
    ID             uuid.UUID       `db:"id" json:"id"`
    UserID         uuid.UUID       `db:"user_id" json:"user_id"`
    Status         string          `db:"status" json:"status"`
    TotalAmount    float64         `db:"total_amount" json:"total_amount"`
    ShippingAddress Address        `db:"shipping_address" json:"shipping_address"`
    PaymentMethod  string          `db:"payment_method" json:"payment_method"`
    CreatedAt      time.Time       `db:"created_at" json:"created_at"`
    UpdatedAt      time.Time       `db:"updated_at" json:"updated_at"`
}
