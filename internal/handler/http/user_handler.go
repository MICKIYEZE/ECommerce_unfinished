package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	userService "ecommerce/internal/service/user"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
    userService userService.UserService
}

func NewUserHandler(svc userService.UserService) *UserHandler {
    return &UserHandler{
        userService: svc,
    }
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
    r.Route("/profile", func(r chi.Router) {
        r.Get("/", h.GetProfile)
        r.Put("/", h.UpdateProfile)
    })
}

// UpdateProfile godoc
// @Summary Update current user's profile
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body UpdateUserRequest true "Updated user data"
// @Success 200 {object} entity.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/me [put]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    userID, ok := GetUserIDFromContext(r.Context())
    if !ok {
        respondError(w, http.StatusUnauthorized, "user not authenticated")
        return
    }

    user, err := h.userService.GetProfile(r.Context(), userID)
    if err != nil {
        respondError(w, http.StatusNotFound, "user not found")
        return
    }

    respondJSON(w, http.StatusOK, user)
}


func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
    userID, ok := GetUserIDFromContext(r.Context())
    if !ok {
        respondError(w, http.StatusUnauthorized, "user not authenticated")
        return
    }

    var req userService.UpdateProfileRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    updatedUser, err := h.userService.UpdateProfile(r.Context(), userID, req, false)
    if err != nil {
        respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, updatedUser)
}

// ListUsers godoc
// @Summary List all users (admin only)
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.User
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
    filter := userService.UserFilter{
        Search: r.URL.Query().Get("search"),
    }

    if limit := r.URL.Query().Get("limit"); limit != "" {
        fmt.Sscanf(limit, "%d", &filter.Limit)
    }
    if offset := r.URL.Query().Get("offset"); offset != "" {
        fmt.Sscanf(offset, "%d", &filter.Offset)
    }

    users, err := h.userService.ListUser(r.Context(), filter)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, users)
}

// DeleteUser godoc
// @Summary Delete a user by ID (admin only)
// @Tags Admin
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    if idParam == "" {
        respondError(w, http.StatusBadRequest, "missing user id")
        return
    }

    userID, err := uuid.Parse(idParam)
    if err != nil {
        respondError(w, http.StatusBadRequest, "invalid user id")
        return
    }

    if err := h.userService.Delete(r.Context(), userID); err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{
        "message": "user deleted",
    })
}



