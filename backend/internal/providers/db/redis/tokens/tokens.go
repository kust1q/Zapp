package tokens

import (
	"github.com/redis/go-redis/v9"
)

const (
	prefixRefreshToken  = "refresh:"
	prefixRecoveryToken = "recovery:"
	prefixUserSessions  = "user_sessions:"
)

type tokenStorage struct {
	redis *redis.Client
}

func NewTokenStorage(redis *redis.Client) *tokenStorage {
	return &tokenStorage{
		redis: redis,
	}
}
