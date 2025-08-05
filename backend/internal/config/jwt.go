package config

import "time"

type JWTConfig struct {
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	PrivateKeyPath string
	PublicKeyPath  string
}
