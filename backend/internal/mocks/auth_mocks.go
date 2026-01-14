package mocks

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type fakeTx struct{}

func (f *fakeTx) Commit() error   { return nil }
func (f *fakeTx) Rollback() error { return nil }

type MockDB struct {
	mock.Mock
}

func (m *MockDB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return nil, nil
}

func (m *MockDB) CreateUserTx(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, tx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockDB) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockDB) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockDB) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockDB) UpdateUserPassword(ctx context.Context, userID int, password string) error {
	args := m.Called(ctx, userID, password)
	return args.Error(0)
}

func (m *MockDB) DeleteUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockDB) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

type MockTokenStorage struct {
	mock.Mock
}

func (m *MockTokenStorage) StoreRefresh(ctx context.Context, refreshToken, userID string, ttl time.Duration) error {
	args := m.Called(ctx, refreshToken, userID, ttl)
	return args.Error(0)
}

func (m *MockTokenStorage) GetUserIdByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockTokenStorage) CloseAllSessions(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockTokenStorage) RemoveRefresh(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockTokenStorage) StoreRecovery(ctx context.Context, token, userID string, ttl time.Duration) error {
	args := m.Called(ctx, token, userID, ttl)
	return args.Error(0)
}

func (m *MockTokenStorage) GetUserIdByRecoveryToken(ctx context.Context, recoveryToken string) (string, error) {
	args := m.Called(ctx, recoveryToken)
	return args.String(0), args.Error(1)
}

func (m *MockTokenStorage) RemoveRecovery(ctx context.Context, recoveryToken string) error {
	args := m.Called(ctx, recoveryToken)
	return args.Error(0)
}

type MockMediaService struct {
	mock.Mock
}

func (m *MockMediaService) UploadAvatarTx(ctx context.Context, userID int, file io.Reader, filename string, tx *sql.Tx) (*entity.Avatar, error) {
	args := m.Called(ctx, userID, file, filename, tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Avatar), args.Error(1)
}

type MockEventProducer struct {
	mock.Mock
}

func (m *MockEventProducer) Publish(ctx context.Context, topic string, event any) error {
	args := m.Called(ctx, topic, event)
	return args.Error(0)
}
