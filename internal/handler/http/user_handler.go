package http

import (
    "encoding/json"
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
