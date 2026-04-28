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

// @Summary Update user UpdateProfile
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /users [put]
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

// Me godoc
// @Summary Get current user
// @Description Returns the authenticated user's profile
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} userService.UserResponse
// @Failure 401 {object} map[string]string
// @Router /users/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
    userID, ok := GetUserIDFromContext(r.Context())
    if !ok {
        respondError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    user, err := h.userService.GetByID(r.Context(), userID)
    if err != nil {
        respondError(w, http.StatusInternalServerError, "Failed to fetch user")
        return
    }

    respondJSON(w, http.StatusOK, user)
}
