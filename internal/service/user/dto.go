package user

import "ecommerce/internal/domain"

type UserFilter struct {
    Search  string  `json:"search"`
    Role    *string `json:"role"`
    Limit   int     `json:"limit"`
    Offset  int     `json:"offset"`
    OrderBy string  `json:"order_by"`
}

type UserListResponse struct {
    Users  []*domain.User `json:"users"`
    Total  int            `json:"total"`
    Limit  int            `json:"limit"`
    Offset int            `json:"offset"`
}

type UpdateProfileRequest struct {
    Email     *string `json:"email"`
    Password  *string `json:"password"`
    FirstName *string `json:"first_name" validate:"omitempty,min=2,max=100"`
    LastName  *string `json:"last_name" validate:"omitempty,min=2,max=100"`
    Role      *string `json:"role"`
}
