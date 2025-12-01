package config

import (
	"crypto/rsa"
	"time"
)

type (
	ApplicationConfig struct {
		Port string
	}

	CacheConfig struct {
		DefaultTtl  time.Duration
		CountersTtl time.Duration
	}

	TokensConfig struct {
		AccessTTL   time.Duration
		RefreshTTL  time.Duration
		RecoveryTTL time.Duration
	}

	JWTConfig struct {
		PrivateKey *rsa.PrivateKey
		PublicKey  *rsa.PublicKey
	}

	MinioConfig struct {
		Port       string
		Endpoint   string
		BucketName string
		User       string
		Password   string
		UseSSL     bool
		TTL        time.Duration
	}

	PostgresConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	RedisConfig struct {
		Host     string
		Port     string
		Password string
		DB       int
	}

	ElasticConfig struct {
		Host string
		Port string
	}
)
