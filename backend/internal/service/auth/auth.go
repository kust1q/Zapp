package auth

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"io"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
)

var (
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrInvalidGender       = errors.New("invalid gender")
	ErrCacheUnavailable    = errors.New("cache service unavailable")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrInvalidSecretAnswer = errors.New("Invalid secret answer")
	ErrInvalidCredentials  = errors.New("invalid credential")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrTokenNotFound       = errors.New("refresh token not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotFound        = errors.New("user not found")
)

type userCache interface {
	Add(ctx context.Context, dataType, data string) error
}

type userStorage interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateUserTx(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error)
	SetSecretQuestionTx(ctx context.Context, tx *sql.Tx, question *entity.SecretQuestion) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	SetSecretQuestion(ctx context.Context, question *entity.SecretQuestion) error
	GetSecurityDataByUserID(ctx context.Context, userID int) (*entity.SecretQuestion, error)
	UpdateUserPassword(ctx context.Context, userID int, password string) error
	DeleteUser(ctx context.Context, userID int) error
	UserExistsByUsername(ctx context.Context, username string) (bool, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
}

type tokenStorage interface {
	Store(ctx context.Context, token, userID string, ttl time.Duration) error
	GetUserIdByRefreshToken(ctx context.Context, token string) (string, error)
	CloseAllSessions(ctx context.Context, userID string) error
	Remove(ctx context.Context, token string) error
}

type mediaService interface {
	UploadAvatarTx(ctx context.Context, userID int, file io.Reader, filename string, tx *sql.Tx) (*entity.Avatar, error)
	GetPresignedURL(ctx context.Context, objectPath string) (string, error)
}

type AuthServiceConfig struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type authService struct {
	cfg     AuthServiceConfig
	storage userStorage
	cache   userCache
	media   mediaService
	tokens  tokenStorage
}

func NewAuthService(cfg AuthServiceConfig, storage userStorage, cache userCache, media mediaService, tokens tokenStorage) *authService {
	return &authService{
		cfg:     cfg,
		storage: storage,
		cache:   cache,
		media:   media,
		tokens:  tokens,
	}
}
