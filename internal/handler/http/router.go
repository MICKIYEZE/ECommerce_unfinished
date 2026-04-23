package http

import (
	"net/http"

	authService "ecommerce/internal/service/auth"
	userService "ecommerce/internal/service/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterConfig struct {
    UserService userService.UserService
    AuthService authService.AuthService
}

func NewRouter(config RouterConfig) *chi.Mux {
    r := chi.NewRouter()

    // Global middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)

    // Health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        respondJSON(w, http.StatusOK, map[string]string{
            "status": "ok",
        })
    })

    // API routes
    r.Route("/api/v1", func(r chi.Router) {
        userHandler := NewUserHandler(config.UserService)
        userHandler.RegisterRoutes(r)
    })

    r.Group(func(r chi.Router) {
        // r.Use(RequireAuth(config.AuthService))
        r.Use(RequireAdmin)
    })

    return r
}
