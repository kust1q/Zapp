package data

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	prefixRefreshToken = "refresh:"
)

type tokenStorage struct {
	redis *redis.Client
}

func NewTokenStorage(redis *redis.Client) *tokenStorage {
	return &tokenStorage{
		redis: redis,
	}
}

func (s *tokenStorage) Store(ctx context.Context, token, userID string, ttl time.Duration) error {
	return s.redis.Set(ctx, s.buildKey(token), userID, ttl).Err()
}

func (s *tokenStorage) Exists(ctx context.Context, token, userID string) (bool, error) {
	exists, err := s.redis.Get(ctx, s.buildKey(token)).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, fmt.Errorf("redis error: %w", err)
	}
	return exists == userID, nil
}

func (s *tokenStorage) Remove(ctx context.Context, token, userID string) error {
	return s.redis.Del(ctx, s.buildKey(token)).Err()
}

func (s *tokenStorage) buildKey(refreshToken string) string {
	return prefixRefreshToken + refreshToken
}
