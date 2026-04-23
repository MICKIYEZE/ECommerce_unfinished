package entity

type ProductFilter struct {
    Search     string  `json:"search"`
    CategoryID *string `json:"category_id"`
    MinPrice   *float64 `json:"min_price"`
    MaxPrice   *float64 `json:"max_price"`
    SortBy     string  `json:"sort_by"`
    Limit      int     `json:"limit"`
    Offset     int     `json:"offset"`
}
