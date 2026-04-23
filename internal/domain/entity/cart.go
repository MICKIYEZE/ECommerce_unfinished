package entity

type Cart struct {
    Items []CartItem `json:"items"`
    Total float64    `json:"total"`
}
