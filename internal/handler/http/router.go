package http

import (
    "net/http"

    authService "ecommerce/internal/service/auth"
    userService "ecommerce/internal/service/user"
    productService "ecommerce/internal/service/product"
    cartService "ecommerce/internal/service/cart"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    httpSwagger "github.com/swaggo/http-swagger"
)

type RouterConfig struct {
    UserService    userService.UserService
    AuthService    authService.AuthService
    ProductService productService.ProductService
    CartService    cartService.CartService
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

        // Public user routes
        userHandler.RegisterRoutes(r)

        // Protected user routes
        r.Group(func(r chi.Router) {
            r.Use(RequireAuth(config.AuthService))
            r.Put("/users/me", userHandler.UpdateProfile)
        })

        // ---------------------------
        // PRODUCT ROUTES
        // ---------------------------
        productHandler := NewProductHandler(config.ProductService)

        r.Route("/products", func(r chi.Router) {
            r.Get("/", productHandler.ListProducts)
            r.Get("/{id}", productHandler.GetProduct)

            r.Group(func(r chi.Router) {
                r.Use(RequireAuth(config.AuthService))
                r.Use(RequireAdmin)

                r.Post("/", productHandler.CreateProduct)
                r.Put("/{id}", productHandler.UpdateProduct)
                r.Delete("/{id}", productHandler.DeleteProduct)
            })
        })

        // ---------------------------
        // CART ROUTES
        // ---------------------------
        cartHandler := NewCartHandler(config.CartService)

        r.Route("/cart", func(r chi.Router) {
            r.Group(func(r chi.Router) {
                // r.Use(RequireAuth(config.AuthService))

                r.Get("/", cartHandler.GetCart)
                r.Post("/add", cartHandler.AddItem)
                r.Put("/update", cartHandler.UpdateItem)
                r.Delete("/remove/{productID}", cartHandler.RemoveItem)
                r.Delete("/clear", cartHandler.ClearCart)
            })
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
