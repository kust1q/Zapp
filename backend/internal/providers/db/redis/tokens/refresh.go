package tokens

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func (s *tokenStorage) StoreRefresh(ctx context.Context, refreshToken, userID string, ttl time.Duration) error {
	refreshKey, err := s.buildRefreshKey(refreshToken)
	if err != nil {
		return err
	}
	_, err = s.redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SAdd(ctx, refreshKey, refreshToken)
		pipe.Expire(ctx, s.buildSessionKey(userID), ttl)
		pipe.Set(ctx, refreshKey, userID, ttl)
		return nil
	})
	return err
}

func (s *tokenStorage) GetUserIdByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	refreshKey, err := s.buildRefreshKey(refreshToken)
	if err != nil {
		return "", err
	}
	userID, err := s.redis.Get(ctx, refreshKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("redis error: %w", err)
	}
	return userID, nil
}

func (s *tokenStorage) RemoveRefresh(ctx context.Context, refreshToken string) error {
	userID, err := s.GetUserIdByRefreshToken(ctx, refreshToken)
	if userID == "" {
		return nil
	}
	if err != nil {
		return err
	}
	refreshKey, err := s.buildRefreshKey(refreshToken)
	if err != nil {
		return err
	}
	pipe := s.redis.Pipeline()
	pipe.Del(ctx, refreshKey)
	pipe.SRem(ctx, s.buildSessionKey(userID), refreshToken)
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
		refreshKey, err := s.buildRefreshKey(token)
		if err != nil {
			return err
		}
		pipe.Del(ctx, refreshKey)
		pipe.SRem(ctx, s.buildSessionKey(userID), token)
	}
	_, err = pipe.Exec(ctx)
	return err
}

func (s *tokenStorage) buildRefreshKey(refreshToken string) (string, error) {
	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash refresh token: %w", err)
	}
	return prefixRefreshToken + string(refreshHash), nil
}

func (s *tokenStorage) buildSessionKey(userID string) string {
	return prefixUserSessions + userID
}
