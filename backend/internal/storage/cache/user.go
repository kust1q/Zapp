package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	EmailType    = "email"
	UsernameType = "username"

	prefixEmail    = "e:"
	prefixUsername = "u:"
)

var (
	userMap = map[string]string{
		UsernameType: prefixUsername,
		EmailType:    prefixEmail,
	}
)

type hasher interface {
	AuthHash(dataType, data string) string
}

type authCache struct {
	redis  *redis.Client
	hasher hasher
	ttl    time.Duration
}

func NewAuthCache(redis *redis.Client, hasher hasher, ttl time.Duration) *authCache {
	return &authCache{
		redis:  redis,
		hasher: hasher,
		ttl:    ttl,
	}
}

func (c *authCache) Add(ctx context.Context, dataType, data string) error {
	prefix, ok := userMap[dataType]
	if !ok {
		return fmt.Errorf("wrong user data type: %s", dataType)
	}

	key := prefix + c.hasher.AuthHash(dataType, data)
	return c.redis.Set(ctx, key, "", c.ttl).Err()
}

func (c *authCache) Exists(ctx context.Context, dataType, data string) (bool, error) {
	key, err := c.buildKey(dataType, data)
	if err != nil {
		return false, err
	}
	exists, err := c.redis.Exists(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, fmt.Errorf("redis error: %w", err)
	}
	return exists == 1, nil
}

func (c *authCache) buildKey(dataType, data string) (string, error) {
	prefix, ok := userMap[dataType]
	if !ok {
		return "", fmt.Errorf("invalid user data type: %s", dataType)
	}
	return prefix + c.hasher.AuthHash(dataType, data), nil
}

/*
func (c *AuthCache) HealthCheck(ctx context.Context) error {
	_, err := c.client.Ping(ctx).Result()
	return err
}
*/
