package tokens

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func (s *tokenStorage) StoreRecovery(ctx context.Context, recoveryToken, userID string, ttl time.Duration) error {
	recoveryKey, err := s.buildRecoveryKey(recoveryToken)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, recoveryKey, userID, ttl).Err()
}

func (s *tokenStorage) GetUserIdByRecoveryToken(ctx context.Context, recoveryToken string) (string, error) {
	recoveryKey, err := s.buildRecoveryKey(recoveryToken)
	if err != nil {
		return "", err
	}
	userID, err := s.redis.Get(ctx, recoveryKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("redis error: %w", err)
	}
	return userID, nil
}

func (s *tokenStorage) RemoveRecovery(ctx context.Context, recoveryToken string) error {
	recoveryKey, err := s.buildRecoveryKey(recoveryToken)
	if err != nil {
		return err
	}
	err = s.redis.Del(ctx, recoveryKey).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to remove recovery token: %w", err)
	}
	return nil
}

func (s *tokenStorage) buildRecoveryKey(recoveryToken string) (string, error) {
	recoveryHash, err := bcrypt.GenerateFromPassword([]byte(recoveryToken), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash refresh token: %w", err)
	}
	return prefixRecoveryToken + string(recoveryHash), nil
}
