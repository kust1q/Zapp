package config

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type config struct {
	App      ApplicationConfig
	Postgres PostgresConfig
	Minio    MinioConfig
}

var (
	instance *config
	once     sync.Once
)

func Get() *config {
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
				Port:          viper.GetString("minio.port"),
				MinioEndpoint: viper.GetString("minio.endpoint"),
				BucketName:    viper.GetString("minio.bucketname"),
				MinioUser:     viper.GetString("minio.user"),
				MinioPassword: viper.GetString("minio.password"),
				MinioUseSSL:   viper.GetBool("minio.sslmode"),
			},
		}
	})
	return instance
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
