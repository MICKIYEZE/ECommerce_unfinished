package user_test

import (
    "context"
    "errors"
    "testing"

    "ecommerce/internal/domain/entity"
    "ecommerce/internal/service/user"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

//
// FULL MOCK — matches postgres.UserRepository EXACTLY
//

type MockUserRepo struct {
    CreateFn                func(ctx context.Context, u *entity.User) error
    GetByIDFn               func(ctx context.Context, id uuid.UUID) (*entity.User, error)
    GetByEmailFn            func(ctx context.Context, email string) (*entity.User, error)
    ListFn                  func(ctx context.Context, limit, offset int) ([]*entity.User, error)
    UpdateFn                func(ctx context.Context, u *entity.User) error
    DeleteFn                func(ctx context.Context, id uuid.UUID) error
    StoreRefreshTokenFn     func(ctx context.Context, userID uuid.UUID, token string) error
    GetUserByRefreshTokenFn func(ctx context.Context, token string) (*entity.User, error)
    DeleteRefreshTokenFn    func(ctx context.Context, token string) error
}

func (m *MockUserRepo) Create(ctx context.Context, u *entity.User) error {
    if m.CreateFn != nil {
        return m.CreateFn(ctx, u)
    }
    return nil
}

func (m *MockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    return m.GetByIDFn(ctx, id)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
    if m.GetByEmailFn != nil {
        return m.GetByEmailFn(ctx, email)
    }
    return nil, nil
}

func (m *MockUserRepo) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
    return m.ListFn(ctx, limit, offset)
}

func (m *MockUserRepo) Update(ctx context.Context, u *entity.User) error {
    return m.UpdateFn(ctx, u)
}

func (m *MockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
    return m.DeleteFn(ctx, id)
}

func (m *MockUserRepo) StoreRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
    if m.StoreRefreshTokenFn != nil {
        return m.StoreRefreshTokenFn(ctx, userID, token)
    }
    return nil
}

func (m *MockUserRepo) GetUserByRefreshToken(ctx context.Context, token string) (*entity.User, error) {
    if m.GetUserByRefreshTokenFn != nil {
        return m.GetUserByRefreshTokenFn(ctx, token)
    }
    return nil, nil
}

func (m *MockUserRepo) DeleteRefreshToken(ctx context.Context, token string) error {
    if m.DeleteRefreshTokenFn != nil {
        return m.DeleteRefreshTokenFn(ctx, token)
    }
    return nil
}

//
// TESTS
//

func TestGetProfile_Success(t *testing.T) {
    uid := uuid.New()

    mockRepo := &MockUserRepo{
        GetByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
            return &entity.User{ID: id, FirstName: "John"}, nil
        },
    }

    svc := user.NewService(mockRepo)

    u, err := svc.GetProfile(context.Background(), uid)

    assert.NoError(t, err)
    assert.Equal(t, uid, u.ID)
}

func TestGetProfile_NotFound(t *testing.T) {
    mockRepo := &MockUserRepo{
        GetByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
            return nil, errors.New("not found")
        },
    }

    svc := user.NewService(mockRepo)

    u, err := svc.GetProfile(context.Background(), uuid.New())

    assert.Nil(t, u)
    assert.EqualError(t, err, user.ErrUserNotFound.Error())
}

func TestDelete_Success(t *testing.T) {
    mockRepo := &MockUserRepo{
        DeleteFn: func(ctx context.Context, id uuid.UUID) error { return nil },
    }

    svc := user.NewService(mockRepo)

    err := svc.Delete(context.Background(), uuid.New())

    assert.NoError(t, err)
}

func TestDelete_Error(t *testing.T) {
    mockRepo := &MockUserRepo{
        DeleteFn: func(ctx context.Context, id uuid.UUID) error {
            return errors.New("db error")
        },
    }

    svc := user.NewService(mockRepo)

    err := svc.Delete(context.Background(), uuid.New())

    assert.EqualError(t, err, "failed to delete user: db error")
}

func TestListUser_Filtering(t *testing.T) {
    mockRepo := &MockUserRepo{
        ListFn: func(ctx context.Context, limit, offset int) ([]*entity.User, error) {
            return []*entity.User{
                {FirstName: "John", LastName: "Doe", Role: "admin"},
                {FirstName: "Jane", LastName: "Smith", Role: "customer"},
            }, nil
        },
    }

    svc := user.NewService(mockRepo)

    role := "admin"
    filter := user.UserFilter{
        Role:   &role,
        Search: "jo",
        Limit:  10,
        Offset: 0,
    }

    resp, err := svc.ListUser(context.Background(), filter)

    assert.NoError(t, err)
    assert.Len(t, resp.Users, 1)
    assert.Equal(t, "John", resp.Users[0].FirstName)
}

func TestUpdateProfile_Success(t *testing.T) {
    uid := uuid.New()

    mockRepo := &MockUserRepo{
        GetByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
            return &entity.User{
                ID:          id,
                FirstName:   "Old",
                LastName:    "Name",
                Email:       "old@mail.com",
                Role:        "customer",
                PasswordHash: "oldhash",
            }, nil
        },
        UpdateFn: func(ctx context.Context, u *entity.User) error { return nil },
    }

    svc := user.NewService(mockRepo)

    newFirst := "New"
    newLast := "User"
    newEmail := "new@mail.com"

    req := user.UpdateProfileRequest{
        FirstName: &newFirst,
        LastName:  &newLast,
        Email:     &newEmail,
    }

    updated, err := svc.UpdateProfile(context.Background(), uid, req, false)

    assert.NoError(t, err)
    assert.Equal(t, "New", updated.FirstName)
    assert.Equal(t, "User", updated.LastName)
    assert.Equal(t, "new@mail.com", updated.Email)
}

func TestUpdateProfile_HashPassword(t *testing.T) {
    uid := uuid.New()

    mockRepo := &MockUserRepo{
        GetByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
            return &entity.User{ID: id}, nil
        },
        UpdateFn: func(ctx context.Context, u *entity.User) error { return nil },
    }

    svc := user.NewService(mockRepo)

    pass := "Secret123"

    req := user.UpdateProfileRequest{
        Password: &pass,
    }

    updated, err := svc.UpdateProfile(context.Background(), uid, req, false)

    assert.NoError(t, err)
    assert.NotEqual(t, pass, updated.PasswordHash)
    assert.True(t, len(updated.PasswordHash) > 20)
}

func TestUpdateProfile_AdminRoleChange(t *testing.T) {
    uid := uuid.New()

    mockRepo := &MockUserRepo{
        GetByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
            return &entity.User{ID: id, Role: "customer"}, nil
        },
        UpdateFn: func(ctx context.Context, u *entity.User) error { return nil },
    }

    svc := user.NewService(mockRepo)

    newRole := "admin"

    req := user.UpdateProfileRequest{
        Role: &newRole,
    }

    updated, err := svc.UpdateProfile(context.Background(), uid, req, true)

    assert.NoError(t, err)
    assert.Equal(t, "admin", updated.Role)
}

func TestUpdateProfile_UserNotFound(t *testing.T) {
    mockRepo := &MockUserRepo{
        GetByIDFn: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
            return nil, errors.New("not found")
        },
    }

    svc := user.NewService(mockRepo)

    _, err := svc.UpdateProfile(context.Background(), uuid.New(), user.UpdateProfileRequest{}, false)

    assert.EqualError(t, err, user.ErrUserNotFound.Error())
}
