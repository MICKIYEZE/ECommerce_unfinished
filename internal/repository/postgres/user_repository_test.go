package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"ecommerce/internal/db"
	"ecommerce/internal/domain/entity"
)

func TestUserRepository_Create(t *testing.T) {
    // Create sqlmock DB
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    user := &entity.User{
        ID:           uuid.New(),
        Email:        "test@example.com",
        PasswordHash: "hash",
        FirstName:    "John",
        LastName:     "Doe",
        Role:         "customer",
    }

    createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

    mock.ExpectQuery(`INSERT INTO users\s*\(.*\)\s*VALUES\s*\(.*\)\s*RETURNING created_at, updated_at`).
        WithArgs(
            sqlmock.AnyArg(), 
            user.Email,
            user.PasswordHash,
            user.FirstName,
            user.LastName,
            user.Role,
        ).
        WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
            AddRow(createdAt, createdAt))

    err := repo.Create(context.Background(), user)

    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    userID := uuid.New()
    createdAt := time.Now().UTC()

    mock.ExpectQuery(`SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at FROM users WHERE id = \$1`).
        WithArgs(userID).
        WillReturnRows(sqlmock.NewRows([]string{
            "id", "email", "password_hash", "first_name", "last_name", "role", "created_at", "updated_at",
        }).AddRow(
            userID,
            "test@example.com",
            "hash",
            "John",
            "Doe",
            "customer",
            createdAt,
            createdAt,
        ))

    user, err := repo.GetByID(context.Background(), userID)

    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, userID, user.ID)
    assert.NoError(t, mock.ExpectationsWereMet())
}


func TestUserRepository_GetByEmail(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    email := "test@example.com"
    createdAt := time.Now().UTC()
    userID := uuid.New()

    mock.ExpectQuery(`SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at FROM users WHERE email = \$1`).
        WithArgs(email).
        WillReturnRows(sqlmock.NewRows([]string{
            "id", "email", "password_hash", "first_name", "last_name", "role", "created_at", "updated_at",
        }).AddRow(
            userID,
            email,
            "hash",
            "John",
            "Doe",
            "customer",
            createdAt,
            createdAt,
        ))

    user, err := repo.GetByEmail(context.Background(), email)

    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, email, user.Email)
    assert.NoError(t, mock.ExpectationsWereMet())
}


func TestUserRepository_Update(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    user := &entity.User{
        ID:        uuid.New(),
        FirstName: "Updated",
        LastName:  "User",
    }

    updatedAt := time.Now().UTC()

    mock.ExpectQuery(`UPDATE users\s+SET\s+first_name = \$1,\s+last_name = \$2,\s+updated_at = NOW\(\)\s+WHERE id = \$3\s+RETURNING updated_at`).
        WithArgs(
            user.FirstName,
            user.LastName,
            user.ID,
        ).
        WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).
            AddRow(updatedAt))

    err := repo.Update(context.Background(), user)

    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    userID := uuid.New()

    mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
        WithArgs(userID).
        WillReturnResult(sqlmock.NewResult(0, 1))

    err := repo.Delete(context.Background(), userID)

    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_NoRows(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    userID := uuid.New()

    mock.ExpectQuery(`SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at FROM users WHERE id = \$1`).
        WithArgs(userID).
        WillReturnError(sql.ErrNoRows)

    user, err := repo.GetByID(context.Background(), userID)

    assert.Error(t, err)
    assert.Nil(t, user)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail_NoRows(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    email := "missing@example.com"

    mock.ExpectQuery(`SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at FROM users WHERE email = \$1`).
        WithArgs(email).
        WillReturnError(sql.ErrNoRows)

    user, err := repo.GetByEmail(context.Background(), email)

    assert.Error(t, err)
    assert.Nil(t, user)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update_SQL_Error(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    user := &entity.User{
        ID:        uuid.New(),
        FirstName: "Bad",
        LastName:  "User",
    }

    mock.ExpectQuery(`UPDATE users`).
        WithArgs(user.FirstName, user.LastName, user.ID).
        WillReturnError(errors.New("update failed"))

    err := repo.Update(context.Background(), user)

    assert.Error(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete_SQL_Error(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    userID := uuid.New()

    mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
        WithArgs(userID).
        WillReturnError(errors.New("delete failed"))

    err := repo.Delete(context.Background(), userID)

    assert.Error(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Create_SQL_Error(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    user := &entity.User{
        ID:           uuid.New(),
        Email:        "fail@example.com",
        PasswordHash: "hash",
        FirstName:    "Fail",
        LastName:     "User",
        Role:         "customer",
    }

    mock.ExpectQuery(`INSERT INTO users`).
        WithArgs(
            sqlmock.AnyArg(),
            user.Email,
            user.PasswordHash,
            user.FirstName,
            user.LastName,
            user.Role,
        ).
        WillReturnError(errors.New("insert failed"))

    err := repo.Create(context.Background(), user)

    assert.Error(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update_NoRowsAffected(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    user := &entity.User{
        ID:        uuid.New(),
        FirstName: "John",
        LastName:  "Doe",
    }

    mock.ExpectQuery(`UPDATE users`).
        WithArgs(user.FirstName, user.LastName, user.ID).
        WillReturnRows(sqlmock.NewRows([]string{"updated_at"})) // no rows returned

    err := repo.Update(context.Background(), user)

    assert.Error(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Delete_NoRowsAffected(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    userID := uuid.New()

    mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
        WithArgs(userID).
        WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

    err := repo.Delete(context.Background(), userID)

    assert.Error(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Create_ScanError(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    user := &entity.User{
        ID:           uuid.New(),
        Email:        "test@example.com",
        PasswordHash: "hash",
        FirstName:    "John",
        LastName:     "Doe",
        Role:         "customer",
    }

    // Return wrong type to force scan failure
    mock.ExpectQuery(`INSERT INTO users`).
        WithArgs(
            sqlmock.AnyArg(),
            user.Email,
            user.PasswordHash,
            user.FirstName,
            user.LastName,
            user.Role,
        ).
        WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
            AddRow("not-a-time", "not-a-time"))

    err := repo.Create(context.Background(), user)

    assert.Error(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByID_ScanError(t *testing.T) {
    sqlDB, mock, _ := sqlmock.New()
    sqlxDB := sqlx.NewDb(sqlDB, "postgres")
    wrappedDB := &db.DB{DB: sqlxDB}

    repo := NewUserRepository(wrappedDB)

    userID := uuid.New()

    mock.ExpectQuery(`SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at FROM users WHERE id = \$1`).
        WithArgs(userID).
        WillReturnRows(sqlmock.NewRows([]string{
            "id", "email", "password_hash", "first_name", "last_name", "role", "created_at", "updated_at",
        }).AddRow(
            "not-a-uuid", // forces scan error
            "email",
            "hash",
            "John",
            "Doe",
            "customer",
            time.Now(),
            time.Now(),
        ))

    user, err := repo.GetByID(context.Background(), userID)

    assert.Error(t, err)
    assert.Nil(t, user)
    assert.NoError(t, mock.ExpectationsWereMet())
}

