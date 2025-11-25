package config

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type config struct {
	App      ApplicationConfig
	Postgres PostgresConfig
	Minio    MinioConfig
	Redis    RedisConfig
	Cache    CacheConfig
	Tokens   TokensConfig
	JWT      JWTConfig
}

var (
	instance *config
	once     sync.Once
)

func Get() *config {
	privateKey, publicKey, err := loadRSAKeys(viper.GetString("jwt.private"), viper.GetString("jwt.public"))
	if err != nil {
		logrus.Fatalf("Failed to load RSA keys: %v", err)
	}

	once.Do(func() {
		instance = &config{
			App: ApplicationConfig{
				Port: viper.GetString("port"),
			},
			Postgres: PostgresConfig{
				Host:     viper.GetString("db.host"),
				Port:     viper.GetString("db.port"),
				User:     viper.GetString("db.user"),
				Password: viper.GetString("db.password"),
				DBName:   viper.GetString("db.name"),
				SSLMode:  viper.GetString("db.sslmode"),
			},
			Minio: MinioConfig{
				Port:       viper.GetString("minio.port"),
				Endpoint:   viper.GetString("minio.endpoint"),
				BucketName: viper.GetString("minio.bucketname"),
				User:       viper.GetString("minio.user"),
				Password:   viper.GetString("minio.password"),
				UseSSL:     viper.GetBool("minio.sslmode"),
				TTL:        viper.GetDuration("minio.ttl"),
			},
			Redis: RedisConfig{
				Host:     viper.GetString("redis.host"),
				Port:     viper.GetString("redis.port"),
				Password: viper.GetString("redis.password"),
				DB:       viper.GetInt("redis.db"),
			},
			Cache: CacheConfig{
				HashSecret: viper.GetString("cache.secret"),
				TTL:        viper.GetDuration("cache.ttl"),
			},
			Tokens: TokensConfig{
				AccessTTL:   viper.GetDuration("tokens.accessTTL"),
				RefreshTTL:  viper.GetDuration("tokens.refreshTTL"),
				RecoveryTTL: viper.GetDuration("tokens.recoveryTTL"),
			},
			JWT: JWTConfig{
				PrivateKey: privateKey,
				PublicKey:  publicKey,
			},
		}
	})
	return instance
}

func (c *config) Validate() error {
	var allErrs []string

	if c.App.Port == "" {
		allErrs = append(allErrs, "app: port is required")
	}

	pg := c.Postgres
	validModes := map[string]bool{
		"disable": true, "allow": true, "prefer": true,
		"require": true, "verify-ca": true, "verify-full": true,
	}
	if pg.Host == "" {
		allErrs = append(allErrs, "postgres: host is required")
	}
	if pg.Port == "" {
		allErrs = append(allErrs, "postgres: port is required")
	}
	if pg.User == "" {
		allErrs = append(allErrs, "postgres: user is required")
	}
	if pg.Password == "" {
		allErrs = append(allErrs, "postgres: password is required")
	}
	if pg.DBName == "" {
		allErrs = append(allErrs, "postgres: dbname is required")
	}
	if !validModes[pg.SSLMode] {
		allErrs = append(allErrs, "postgres: invalid sslmode")
	}

	minio := c.Minio
	if minio.Port == "" {
		allErrs = append(allErrs, "minio: port is required")
	}
	if minio.Endpoint == "" {
		allErrs = append(allErrs, "minio: endpoint is required")
	}
	if minio.BucketName == "" {
		allErrs = append(allErrs, "minio: bucketname is required")
	}
	if minio.User == "" {
		allErrs = append(allErrs, "minio: user is required")
	}
	if minio.Password == "" {
		allErrs = append(allErrs, "minio: password is required")
	}
	if minio.TTL <= 0 {
		allErrs = append(allErrs, "minio: ttl is required")
	}

	redis := c.Redis
	if redis.Host == "" {
		allErrs = append(allErrs, "redis: host is required")
	}
	if redis.Port == "" {
		allErrs = append(allErrs, "redis: port is required")
	}
	if redis.DB < 0 {
		allErrs = append(allErrs, "redis: db must be >= 0")
	}

	if c.Cache.HashSecret == "" {
		allErrs = append(allErrs, "hash: secret is required")
	}
	if c.Cache.TTL <= 0 {
		allErrs = append(allErrs, "hash: ttl must be > 0")
	}

	if c.Tokens.AccessTTL <= 0 {
		allErrs = append(allErrs, "tokens: access ttl must be > 0")
	}
	if c.Tokens.RefreshTTL <= 0 {
		allErrs = append(allErrs, "tokens: refresh ttl must be > 0")
	}
	if c.Tokens.RecoveryTTL <= 0 {
		allErrs = append(allErrs, "tokens: recovery ttl must be > 0")
	}
	if c.JWT.PrivateKey == nil {
		allErrs = append(allErrs, "jwt: private key path is required")
	}
	if c.JWT.PublicKey == nil {
		allErrs = append(allErrs, "jwt: public key path is required")
	}
	if len(allErrs) > 0 {
		return errors.New("config validation errors:\n  • " + strings.Join(allErrs, "\n  • "))
	}
	return nil
}

func InitConfig() error {
	viper.AutomaticEnv()

	viper.SetConfigFile(".env")
	if err := viper.MergeInConfig(); err != nil {
		logrus.Printf("warning: .env file not loaded: %v", err)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("../../configs")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil

}

func loadRSAKeys(privateKeyPath, publicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privatePEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read public key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return privateKey, publicKey, nil
}
