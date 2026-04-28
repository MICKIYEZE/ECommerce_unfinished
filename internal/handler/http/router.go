package http

import (
    "net/http"

    authService "ecommerce/internal/service/auth"
    userService "ecommerce/internal/service/user"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    httpSwagger "github.com/swaggo/http-swagger"
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
    r.Use(CORS)

    // Health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        respondJSON(w, http.StatusOK, map[string]string{
            "status": "ok",
        })
    })

    // Swagger
    r.Get("/swagger/*", httpSwagger.WrapHandler)

    // API routes
    r.Route("/api/v1", func(r chi.Router) {

        // ---------------------------
        // AUTH ROUTES (PUBLIC)
        // ---------------------------
        authHandler := NewAuthHandler(config.AuthService)

        r.Post("/auth/register", authHandler.Register)
        r.Post("/auth/login", authHandler.Login)
        r.Post("/auth/refresh", authHandler.Refresh)

        // ---------------------------
        // USER ROUTES
        // ---------------------------
        userHandler := NewUserHandler(config.UserService)

        // Public user routes (if any)
        userHandler.RegisterRoutes(r)

        // Protected user routes
        r.Group(func(r chi.Router) {
            r.Use(RequireAuth(config.AuthService))

            r.Put("/users/me", userHandler.UpdateProfile)
        })

        // ---------------------------
        // ADMIN ROUTES
        // ---------------------------
        r.Group(func(r chi.Router) {
            r.Use(RequireAuth(config.AuthService))
            r.Use(RequireAdmin)

            r.Get("/admin/users", userHandler.ListUsers)
            r.Delete("/admin/users/{id}", userHandler.DeleteUser)
        })
    })

    return r
}
