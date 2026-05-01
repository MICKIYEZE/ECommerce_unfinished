package http

import (
    "context"
    "net/http"
    "strings"

    authService "ecommerce/internal/service/auth"
    "github.com/google/uuid"
)

type contextKey string

const (
    ContextKeyUserID    contextKey = "user_id"
    ContextKeyUserRole  contextKey = "user_role"
    ContextKeyUserEmail contextKey = "user_email"
)

// func RequireAuth(_ authService.AuthService) func(http.Handler) http.Handler {
//     return func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

//             // TEMPORARY: fake user ID so protected routes work
//             fakeUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

//             ctx := context.WithValue(r.Context(), "userID", fakeUserID)
//             next.ServeHTTP(w, r.WithContext(ctx))
//         })
//     }
// }



func RequireAuth(authSrv authService.AuthService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                respondError(w, http.StatusUnauthorized, "missing auth header")
                return
            }

            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                respondError(w, http.StatusUnauthorized, "invalid auth header format")
                return
            }

            tokenString := parts[1]

            claims, err := authSrv.ValidateToken(tokenString)
            if err != nil {
                respondError(w, http.StatusUnauthorized, "invalid or expired token")
                return
            }

            ctx := r.Context()
            ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
            ctx = context.WithValue(ctx, ContextKeyUserRole, claims.Role)
            ctx = context.WithValue(ctx, ContextKeyUserEmail, claims.Email)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// RequireAdmin ensures the user has admin role
func RequireAdmin(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        role, ok := GetUserRoleFromContext(r.Context())
        if !ok {
            respondError(w, http.StatusUnauthorized, "user not authenticated")
            return
        }

        if role != "admin" {
            respondError(w, http.StatusForbidden, "admin access required")
            return
        }

        next.ServeHTTP(w, r)
    })
}

// CORS middleware
func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
