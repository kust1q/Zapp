package config

type MinioConfig struct {
	Port          string
	MinioEndpoint string
	BucketName    string
	MinioUser     string
	MinioPassword string
	MinioUseSSL   bool
}
