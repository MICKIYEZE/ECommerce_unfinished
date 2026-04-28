package http

import (
    "encoding/json"
    "net/http"

    authService "ecommerce/internal/service/auth"

    _ "ecommerce/internal/service/auth"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
    authService authService.AuthService
}

func NewAuthHandler(authService authService.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}


func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req authService.RegisterRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    user, err := h.authService.Register(r.Context(), req)
    if err != nil {
        handleAuthError(w, err)
        return
    }

    respondJSON(w, http.StatusCreated, user)
}


func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req authService.LoginRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    accessToken, refreshToken, err := h.authService.Login(r.Context(), req)
    if err != nil {
        handleAuthError(w, err)
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}


func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
    var body struct {
        RefreshToken string `json:"refresh_token"`
    }

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        respondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    accessToken, newRefreshToken, err := h.authService.RefreshToken(r.Context(), body.RefreshToken)
    if err != nil {
        handleAuthError(w, err)
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{
        "access_token":  accessToken,
        "refresh_token": newRefreshToken,
    })
}

func handleAuthError(w http.ResponseWriter, err error) {
    switch err.Error() {
    case "invalid credentials":
        respondError(w, http.StatusUnauthorized, "Invalid email or password")
    case "invalid refresh token":
        respondError(w, http.StatusUnauthorized, "Invalid refresh token")
    default:
        respondError(w, http.StatusInternalServerError, "Internal server error")
    }
}
