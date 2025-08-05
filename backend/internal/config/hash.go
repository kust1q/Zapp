package config

import "time"

type CacheConfig struct {
	HashSecret string
	TTL        time.Duration
}
