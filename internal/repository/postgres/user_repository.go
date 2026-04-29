package postgres

import (
    "context"
    "fmt"

    "ecommerce/internal/db"
    "ecommerce/internal/domain/entity"

    "github.com/google/uuid"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, limit, offset int) ([]*entity.User, error)

    StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string) error
    GetUserByRefreshToken(ctx context.Context, token string) (*entity.User, error)
    DeleteRefreshToken(ctx context.Context, token string) error
}

type userRepo struct {
    db *db.DB
}

func NewUserRepository(db *db.DB) UserRepository {
    return &userRepo{db: db}
}

func (u *userRepo) Create(ctx context.Context, user *entity.User) error {
    query := `
        INSERT INTO users (
            id,
            email,
            password_hash,
            first_name,
            last_name,
            role,
            created_at,
            updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        RETURNING created_at, updated_at
    `

    user.ID = uuid.New()

    return u.db.QueryRowContext(
        ctx,
        query,
        user.ID,
        user.Email,
        user.PasswordHash,
        user.FirstName,
        user.LastName,
        user.Role,
    ).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (u *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
    query := `DELETE FROM users WHERE id = $1`

    result, err := u.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("user not found")
    }

    return nil
}

func (u *userRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
    var user entity.User

    query := `
        SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    err := u.db.GetContext(ctx, &user, query, email)
    if err != nil {
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }

    return &user, nil
}

func (u *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    var user entity.User

    query := `
        SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    err := u.db.GetContext(ctx, &user, query, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }

    return &user, nil
}

func (u *userRepo) Update(ctx context.Context, user *entity.User) error {
    query := `
        UPDATE users
        SET first_name = $1,
            last_name = $2,
            updated_at = NOW()
        WHERE id = $3
        RETURNING updated_at
    `

    return u.db.QueryRowContext(
        ctx,
        query,
        user.FirstName,
        user.LastName,
        user.ID,
    ).Scan(&user.UpdatedAt)
}

func (u *userRepo) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
    query := `
        SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

    var users []*entity.User

    err := u.db.SelectContext(ctx, &users, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }

    return users, nil
}

func (u *userRepo) StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
    query := `
        INSERT INTO refresh_tokens (user_id, token)
        VALUES ($1, $2)
    `
    _, err := u.db.ExecContext(ctx, query, userID, token)
    if err != nil {
        return fmt.Errorf("failed to store refresh token: %w", err)
    }
    return nil
}

func (u *userRepo) GetUserByRefreshToken(ctx context.Context, token string) (*entity.User, error) {
    var user entity.User

    query := `
        SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.role, u.created_at, u.updated_at
        FROM users u
        JOIN refresh_tokens rt ON rt.user_id = u.id
        WHERE rt.token = $1
    `

    err := u.db.GetContext(ctx, &user, query, token)
    if err != nil {
        return nil, fmt.Errorf("failed to get user by refresh token: %w", err)
    }

    return &user, nil
}

func (u *userRepo) DeleteRefreshToken(ctx context.Context, token string) error {
    query := `DELETE FROM refresh_tokens WHERE token = $1`

    _, err := u.db.ExecContext(ctx, query, token)
    if err != nil {
        return fmt.Errorf("failed to delete refresh token: %w", err)
    }

    return nil
}
