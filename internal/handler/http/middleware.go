package http

import (
	"context"
	authService "ecommerce/internal/service/auth"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	ContextKeyUserID    contextKey = "user_id"
	ContextKeyUserRole  contextKey = "user_role"
	ContextKeyUserEmail contextKey = "user_email"
)

func RequireAuth(authSrv authService.AuthService) func(http.Handler) http.Handler{
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request))
			authHeader := r.Header.Get("Authorization")
			if authHeader == ""{
				respondError(w, http.StatusUnauthorized, "missing auth header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2|| parts [0] != "Bearer"{
				respondError(w, http.StatusUnauthorized,"Invalid auth header format")
				return
			}

			tokenString := parts[1]

			claims, err := authSrv.ValidateToken(tokenString)
			if err != nil {
				respondError(w,http.StatusUnauthorized,"Invalid or expired token")
				return
			}

			ctx := r.Context()

			ctx = context.WithValue(ctx,ContextKeyUserID,claims.UserID)
			ctx = context.WithValue(ctx,ContextKeyUserRole,claims.Role)
			ctx = context.WithValue(ctx,ContextKeyUserEmail,claims.Email)

 	}
}


func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := GetUserRoleFromContext(r.Context())

		if !ok {
			respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		if role != "admin" {
			respondError(w, http.StatusForbidden, "Admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
// Cors - Cross-Origin Resourse Sharing
func CORS(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(ContextKeyUserID).(uuid.UUID)
	return userID, ok
}

func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextKeyUserRole).(string)
	return role, ok
}

func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(ContextKeyUserEmail).(string)
	return email, ok
}