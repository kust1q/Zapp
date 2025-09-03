package auth

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"io"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	media "github.com/kust1q/Zapp/backend/internal/storage/objects"
)

type userCache interface {
	Add(ctx context.Context, dataType, data string) error
}

type userStorage interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateUserTx(ctx context.Context, tx *sql.Tx, user *entity.User) (entity.User, error)
	SetSecretQuestionTx(ctx context.Context, tx *sql.Tx, question *entity.SecretQuestion) error
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByID(ctx context.Context, userID int) (entity.User, error)
	SetSecretQuestion(ctx context.Context, question *entity.SecretQuestion) error
	GetSecurityDataByUserID(ctx context.Context, userID int) (entity.SecretQuestion, error)
	UpdateUserPassword(ctx context.Context, userID int, password string) error
	DeleteUser(ctx context.Context, userID int) error
	UserExistsByUsername(ctx context.Context, username string) (bool, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
}

type mediaStorage interface {
	Upload(ctx context.Context, file io.Reader, mediaType media.MediaType, filename string) (string, string, error)
	Remove(ctx context.Context, objectPath string) error
}

type tokenStorage interface {
	Store(ctx context.Context, token, userID string, ttl time.Duration) error
	GetUserIdByRefreshToken(ctx context.Context, token string) (string, error)
	CloseAllSessions(ctx context.Context, userID string) error
	Remove(ctx context.Context, token string) error
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
	media   mediaStorage
	tokens  tokenStorage
}

func NewAuthService(cfg AuthServiceConfig, storage userStorage, cache userCache, media mediaStorage, tokens tokenStorage) *authService {
	return &authService{
		cfg:     cfg,
		storage: storage,
		cache:   cache,
		media:   media,
		tokens:  tokens,
	}
}
