package http

import (
	userService "ecommerce/internal/service/user"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type RouterConfig struct{
	//AuthService authService.AuthService
	UserService userService.UserService
}

func newRouter(config RouterConfig) *chi.Mux{
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/health",func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	r.Route("api/v1", func(r chi.Router){
		r.Group(func(r chi.Router) {
			userHandler := NewUserHandler(config.UserService)
			userHandler.RegisterRoutes(r)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(RequireAdmin)
	})

	return r
}