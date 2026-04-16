package service

import models "ecommerce/internal/domain"

type UserFilter struct {
	Search  string
	Role    *string
	Limit   int
	Offset  int
	OrderBy string
}

type UserListResponse struct {
	Users  []*models.User
	Total  int
	Limit  int
	Offset int
}

type UpdateProfileRequest struct {
	Email     *string
	Password  *string
	FirstName *string `json:"first_name" validate:"omniempty, min=2, max=100"`
	LastName  *string `json:"last_name" validate:"omniempty, min=2, max=100"`
	Role      *string
}