package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
    
    user, err := s.userRepo.GetUserByRefreshToken(context.Background(), refreshToken)
    if err != nil || user == nil {
        return "", "", errors.New("invalid refresh token")
    }

    // Delete old token
    _ = s.userRepo.DeleteRefreshToken(context.Background(), refreshToken)

    // Generate new tokens
    accessToken, err := s.generateToken(user)
    if err != nil {
        return "", "", err
    }

    newRefresh := uuid.New().String()
    if err := s.userRepo.StoreRefreshToken(context.Background(), user.ID, newRefresh); err != nil {
        return "", "", err
    }

    return accessToken, newRefresh, nil
}
