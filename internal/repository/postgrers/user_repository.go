package postgres

import (
    "context"
    "fmt"

    "ecommerce/internal/db"
    "ecommerce/internal/domain"

    "github.com/google/uuid"
)

type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
    GetByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

type userRepo struct {
    db *db.DB
}

func NewUserRepository(db *db.DB) UserRepository {
    return &userRepo{db: db}
}

func (u *userRepo) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (id, email, password_hash, first_name, surname, role)
        VALUES ($1, $2, $3, $4, $5, $6)
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
        user.Surname,
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

func (u *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User

    query := `
        SELECT id, email, password_hash, first_name, surname, role, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    err := u.db.GetContext(ctx, &user, query, email)
    if err != nil {
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }

    return &user, nil
}

func (u *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
    var user domain.User

    query := `
        SELECT id, email, password_hash, first_name, surname, role, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    err := u.db.GetContext(ctx, &user, query, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }

    return &user, nil
}

func (u *userRepo) Update(ctx context.Context, user *domain.User) error {
    query := `
        UPDATE users
        SET first_name = $1,
            surname = $2,
            updated_at = NOW()
        WHERE id = $3
        RETURNING updated_at
    `

    return u.db.QueryRowContext(
        ctx,
        query,
        user.FirstName,
        user.Surname,
        user.ID,
    ).Scan(&user.UpdatedAt)
}

func (u *userRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
    query := `
        SELECT id, email, password_hash, first_name, surname, role, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

    var users []*domain.User

    err := u.db.SelectContext(ctx, &users, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }

    return users, nil
}
