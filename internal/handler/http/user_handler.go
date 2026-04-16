package http

import (
    userService "ecommerce/internal/service/user"
    "net/http"

    "github.com/go-chi/chi"
)

type UserHandler struct {
    userService userService.UserService
}

func NewUserHandler(srvc userService.UserService) *UserHandler {
    return &UserHandler{
        userService: srvc,
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
		respondError(w, http.StatusUnauthorized,"user not authenticated")
		return
	}

	var req userService.UpdateProfileRequest

	updReq := userService.UpdateProfileRequest{
		Email: req.Email,
		FirstName: req.FirstName,
		LastName: req.LastName,
		Password: req.Password,
		Role: req.Role,
	}
	
	updUser, err := h.userService.UpdateProfile(r.Context(), userID, updReq, false)
	if err := json.NewDecoder(r.Body).Decode(&req): err != nil{
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to update profile")
		return
	}


	respondJSON(w, http.StatusOK, user)
}
