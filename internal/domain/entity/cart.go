package entity

import "github.com/google/uuid"

type Cart struct {
    UserID uuid.UUID   `json:"user_id"`
    Items  []CartItem  `json:"items"`
    Total  float64     `json:"total"`
}
