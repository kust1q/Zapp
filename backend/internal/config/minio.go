package config

import "time"

type MinioConfig struct {
	Port       string
	Endpoint   string
	BucketName string
	User       string
	Password   string
	UseSSL     bool
	TTL        time.Duration
}
