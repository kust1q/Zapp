package data

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	prefixRefreshToken = "refresh:"
	prefixUserSessions = "user_sessions:"
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
	s.removeExpiredTokens(context.Background(), userID)
	_, err := s.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SAdd(ctx, s.buildSessionKey(userID), token)
		pipe.Expire(ctx, s.buildSessionKey(userID), ttl)
		pipe.Set(ctx, s.buildTokenKey(token), userID, ttl)
		return nil
	})
	return err
}

func (s *tokenStorage) GetUserIdByRefreshToken(ctx context.Context, token string) (string, error) {
	userID, err := s.redis.Get(ctx, s.buildTokenKey(token)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("redis error: %w", err)
	}
	return userID, nil
}

func (s *tokenStorage) Remove(ctx context.Context, token string) error {
	userID, err := s.GetUserIdByRefreshToken(ctx, token)
	if userID == "" {
		return nil
	}
	if err != nil {
		return err
	}
	pipe := s.redis.Pipeline()
	pipe.Del(ctx, s.buildTokenKey(token))
	pipe.SRem(ctx, s.buildSessionKey(userID), token)
	_, err = pipe.Exec(ctx)
	return err
}

func (s *tokenStorage) CloseAllSessions(ctx context.Context, userID string) error {
	tokens, err := s.redis.SMembers(ctx, s.buildSessionKey(userID)).Result()
	if err != nil {
		return err
	}
	pipe := s.redis.Pipeline()
	for _, token := range tokens {
		pipe.Del(ctx, s.buildTokenKey(token))
		pipe.SRem(ctx, s.buildSessionKey(userID), token)
	}
	_, err = pipe.Exec(ctx)
	return err
}

func (s *tokenStorage) removeExpiredTokens(ctx context.Context, userID string) error {
	tokens, err := s.redis.SMembers(ctx, s.buildSessionKey(userID)).Result()
	if err != nil {
		return err
	}
	pipe := s.redis.Pipeline()
	for _, token := range tokens {
		if ttl, _ := s.redis.TTL(ctx, s.buildTokenKey(token)).Result(); ttl < 0 {
			pipe.SRem(ctx, s.buildSessionKey(userID), token)
		}
	}
	_, err = pipe.Exec(ctx)
	return err
}

func (s *tokenStorage) buildTokenKey(refreshToken string) string {
	return prefixRefreshToken + refreshToken
}

func (s *tokenStorage) buildSessionKey(userID string) string {
	return prefixUserSessions + userID
}
