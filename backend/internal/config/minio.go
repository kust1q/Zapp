package config

type MinioConfig struct {
	Port       string
	Endpoint   string
	BucketName string
	User       string
	Password   string
	UseSSL     bool
}
