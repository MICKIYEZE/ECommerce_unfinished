package auth

// RegisterRequest is the payload for user registration
type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// LoginRequest is the payload for user login
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// AuthResponse is returned after login/register/refresh
type AuthResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}
